package productbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	UpdateProduct(ctx context.Context, cond map[string]interface{}, data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error)
}

type UpdateProductCache interface {
	DeleteProduct(ctx context.Context, morekeys ...string) error
}

type updateProductBusiness struct {
	store UpdateProductStorage
	cache UpdateProductCache
}

func NewUpdatePruductBiz(store UpdateProductStorage, cache UpdateProductCache) *updateProductBusiness {
	return &updateProductBusiness{store, cache}
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

	var pagging paggingcommon.Paging
	pagging.Process()

	key := cacheutil.GenerateKey(productmodel.EntityName, cond, pagging, commonfilter.Filter{})

	if err := biz.cache.DeleteProduct(ctx, key); err != nil {
		return nil, common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	return updatedRecord, nil
}
