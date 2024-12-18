package productbusiness

import (
	"context"
	"log"
	"sync"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	UpdateProduct(ctx context.Context, cond map[string]interface{}, data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error)
}

type UpdateProductCache interface {
	DeleteProduct(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type UpdateImageCloud interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type updateProductBusiness struct {
	store UpdateProductStorage
	cache UpdateProductCache
	cloud UpdateImageCloud
}

func NewUpdatePruductBiz(store UpdateProductStorage, cache UpdateProductCache, cloud UpdateImageCloud) *updateProductBusiness {
	return &updateProductBusiness{store, cache, cloud}
}

func (biz *updateProductBusiness) UpdateProduct(ctx context.Context,
	cond map[string]interface{}, data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error) {

	record, err := biz.store.GetProduct(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrCannotGetEntity(productmodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateProduct(ctx, map[string]interface{}{"product_id": record.ProductID}, data, morekeys...)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	// Only proceed with image updates if data.Images is provided and not empty
	if data.Images != nil && len(data.Images) > 0 {
		// 3. Get existing images associated with the product
		imageCond := map[string]interface{}{
			"product_id": record.ProductID,
		}

		oldImages, err := biz.cloud.ListItem(ctx, imageCond)

		if err != nil {
			return nil, common.ErrCannotListEntity(imagemodel.EntityName, err)
		}

		oldImageIDs := make(map[uint64]bool)
		for _, img := range oldImages {
			oldImageIDs[img.ImageID] = true
		}

		// 4. Determine images to add and remove
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

		// 5. Update ProductID for images to add and remove using goroutines
		var wg sync.WaitGroup
		var updateErr error
		var mu sync.Mutex

		// Limit the number of concurrent goroutines to prevent system overload
		sem := make(chan struct{}, 10) // Limit to 10 goroutines

		// Add new images (Update `ProductID` for these images)
		log.Print(imagesToAdd)
		for _, imgID := range imagesToAdd {
			wg.Add(1)
			sem <- struct{}{} // Acquire a token
			go func(id uint64) {
				defer wg.Done()
				defer func() { <-sem }() // Release the token

				err := biz.cloud.UpdateImage(ctx, map[string]interface{}{"image_id": id}, &imagemodel.UpdateImage{
					ProductID: &record.ProductID,
				})
				if err != nil {
					mu.Lock()
					updateErr = err
					mu.Unlock()
				}
			}(imgID)
		}

		// Remove old images (Set `ProductID` to NULL)
		for _, imgID := range imagesToRemove {
			wg.Add(1)
			sem <- struct{}{}
			go func(id uint64) {
				defer wg.Done()
				defer func() { <-sem }()

				err := biz.cloud.DeleteImage(ctx, map[string]interface{}{"image_id": id})
				if err != nil {
					mu.Lock()
					updateErr = err
					mu.Unlock()
				}
			}(imgID)
		}

		wg.Wait()

		if updateErr != nil {
			// Handle error (rollback transaction if needed)
			return nil, common.ErrCannotUpdateEntity("Image", updateErr)
		}
	}

	// 6. Delete product cache
	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(productmodel.EntityName, err)
	}

	if err := biz.cache.DeleteProduct(ctx, key); err != nil {
		return nil, common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, productmodel.EntityName); err != nil {
		return nil, common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	return updatedRecord, nil
}
