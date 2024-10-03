package productbusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	DeleteProduct(ctx context.Context, cond map[string]interface{}, morekyes ...string) error
}

type DeleteProductCache interface {
	DeleteProduct(ctx context.Context, morekeys ...string) error
}

type DeleteImageCloud interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type deleteProductBusiness struct {
	store DeleteProductStorage
	cache DeleteProductCache
	cloud DeleteImageCloud
}

func NewDeleteProductBiz(store DeleteProductStorage, cache DeleteProductCache, cloud DeleteImageCloud) *deleteProductBusiness {
	return &deleteProductBusiness{store, cache, cloud}
}

func (biz *deleteProductBusiness) DeleteProduct(ctx context.Context, cond map[string]interface{}, morekyes ...string) error {

	record, err := biz.store.GetProduct(ctx, cond)

	if err != nil {
		return common.ErrCannotGetEntity(productmodel.EntityName, err)
	}

	if record == nil {
		return common.ErrNotFoundEntity(productmodel.EntityName, err)
	}

	if err := biz.store.DeleteProduct(ctx, cond, morekyes...); err != nil {
		return common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	//oldImageId := record.ImageID
	//if err := biz.cloud.DeleteImage(ctx, map[string]interface{}{"image_id": oldImageId}); err != nil {
	//	log.Print("could not delete image in database, %e", err)
	//	return nil
	//}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"product_id": record.ProductID,
	}

	oldImages, err := biz.cloud.ListItem(ctx, imageCond)

	if err != nil {
		return common.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	//oldImageIDs := make(map[uint64]bool)
	//for _, img := range oldImages {
	//	oldImageIDs[img.ImageID] = true
	//}

	var imagesToRemove []uint64
	for _, imgID := range oldImages {
		imagesToRemove = append(imagesToRemove, imgID.ImageID)
	}

	// 5. Cập nhật `ProductID` cho hình ảnh cần thêm và xóa sử dụng goroutines
	var wg sync.WaitGroup
	var deleteErr error
	var mu sync.Mutex

	// Giới hạn số lượng goroutines chạy đồng thời để tránh làm quá tải hệ thống
	sem := make(chan struct{}, 10) // Giới hạn 10 goroutines

	// Xóa liên kết hình ảnh cũ (Đặt `ProductID` về NULL)
	for _, imgID := range imagesToRemove {
		wg.Add(1)
		sem <- struct{}{}
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }()
			err := biz.cloud.DeleteImage(ctx, map[string]interface{}{"image_id": id})
			if err != nil {
				mu.Lock()
				deleteErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	wg.Wait()

	if deleteErr != nil {
		// Xử lý lỗi (rollback transaction nếu cần)
		return common.ErrCannotUpdateEntity("Image", deleteErr)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekyes,
	})
	if err != nil {
		return common.ErrCannotGenerateKey(productmodel.EntityName, err)
	}

	if err := biz.cache.DeleteProduct(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	return nil
}
