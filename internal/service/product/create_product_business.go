package productbusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

type CreateProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error)
}

type UpdateImage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type createProductBusiness struct {
	store      CreateProductStorage
	imageStore UpdateImage
}

func NewCreateProductBusiness(store CreateProductStorage, imageStore UpdateImage) *createProductBusiness {
	return &createProductBusiness{store, imageStore}
}

func (biz *createProductBusiness) CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error) {
	//data.ImageID = data.Imag
	recordId, err := biz.store.CreateProduct(ctx, data)

	if err != nil {
		return 0, common.ErrCannotCreateEntity(productmodel.EntityName, err)
	}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"product_id": recordId,
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
				ProductID: &recordId,
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

	return recordId, nil
}
