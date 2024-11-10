package categorybusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type CreateCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	CreateCategory(ctx context.Context, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error)
}

type createCategoryBusiness struct {
	store      CreateCategoryStorage
	imageStore UpdateImage
}

func NewCreateCategoryBusiness(store CreateCategoryStorage, imageStore UpdateImage) *createCategoryBusiness {
	return &createCategoryBusiness{store, imageStore}
}

func (biz *createCategoryBusiness) CreateCategory(ctx context.Context, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error) {

	recordID, err := biz.store.CreateCategory(ctx, data)

	if err != nil {
		// Check for MySQL duplicate entry error
		return 0, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"category_id": recordID,
	}

	oldImages, err := biz.imageStore.ListItem(ctx, imageCond)

	if err != nil {
		return 0, common.ErrCannotListEntity(imagemodel.EntityName, err)
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

			err := biz.imageStore.UpdateImage(ctx, map[string]interface{}{"image_id": id}, &imagemodel.UpdateImage{
				ProductID: &recordID,
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

			err := biz.imageStore.DeleteImage(ctx, map[string]interface{}{"image_id": id})
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
		return 0, common.ErrCannotUpdateEntity("Image", updateErr)
	}

	return recordID, nil
}
