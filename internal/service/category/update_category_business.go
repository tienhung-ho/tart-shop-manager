package categorybusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error)
}

type UpdateCategoryCache interface {
	DeleteCategory(ctx context.Context, morekeys ...string) error
}

type UpdateImage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type updateCategoryBusiness struct {
	store UpdateCategoryStorage
	cache UpdateCategoryCache
	image UpdateImage
}

func NewUpdateCategoryBiz(store UpdateCategoryStorage, cache UpdateCategoryCache, image UpdateImage) *updateCategoryBusiness {
	return &updateCategoryBusiness{store, cache, image}
}

func (biz *updateCategoryBusiness) UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error) {

	record, err := biz.store.GetCategory(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateCategory(ctx, map[string]interface{}{"category_id": record.CategoryID}, data)
	if err != nil {
		return nil, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"category_id": record.CategoryID,
	}

	oldImages, err := biz.image.ListItem(ctx, imageCond)

	if err != nil {
		return nil, common.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	oldImageIDs := make(map[uint64]bool)
	for _, img := range oldImages {
		oldImageIDs[img.ImageID] = true
	}

	// 4. Xác định hình ảnh cần thêm và cần xóa
	newImageIDs := make(map[uint64]bool)
	for _, imgID := range data.Images {
		newImageIDs[imgID.ImageID] = true
	}

	var imagesToAdd []uint64
	var imagesToRemove []uint64

	for imgID := range newImageIDs {
		if !oldImageIDs[imgID] {
			imagesToAdd = append(imagesToAdd, imgID)
		}
	}

	for imgID := range oldImageIDs {
		if !newImageIDs[imgID] {
			imagesToRemove = append(imagesToRemove, imgID)
		}
	}

	// 5. Cập nhật `ProductID` cho hình ảnh cần thêm và xóa sử dụng goroutines
	var wg sync.WaitGroup
	var updateErr error
	var mu sync.Mutex

	// Giới hạn số lượng goroutines chạy đồng thời để tránh làm quá tải hệ thống
	sem := make(chan struct{}, 10) // Giới hạn 10 goroutines

	// Thêm hình ảnh mới (Cập nhật `ProductID` cho các hình ảnh này)
	for _, imgID := range imagesToAdd {
		wg.Add(1)
		sem <- struct{}{} // Acquire a token
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			err := biz.image.UpdateImage(ctx, map[string]interface{}{"image_id": id}, &imagemodel.UpdateImage{
				CategoryID: &record.CategoryID,
			})
			if err != nil {
				mu.Lock()
				updateErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	// Xóa liên kết hình ảnh cũ (Đặt `ProductID` về NULL)
	for _, imgID := range imagesToRemove {
		wg.Add(1)
		sem <- struct{}{}
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }()

			err := biz.image.DeleteImage(ctx, map[string]interface{}{"image_id": id})
			if err != nil {
				mu.Lock()
				updateErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	wg.Wait()

	if updateErr != nil {
		// Xử lý lỗi (rollback transaction nếu cần)
		return nil, common.ErrCannotUpdateEntity("Image", updateErr)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: categorymodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(categorymodel.EntityName, err)
	}

	if err := biz.cache.DeleteCategory(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	return updatedRecord, nil
}
