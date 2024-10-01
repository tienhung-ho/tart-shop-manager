package productbusiness

import (
	"context"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
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

	oldImageId := record.ImageID
	if err := biz.cloud.DeleteImage(ctx, map[string]interface{}{"image_id": oldImageId}); err != nil {
		log.Print("could not delete image in database, %e", err)
		return nil
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
