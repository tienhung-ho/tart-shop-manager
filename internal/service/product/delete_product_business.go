package productbusiness

import (
	"context"
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

type deleteProductBusiness struct {
	store DeleteProductStorage
	cache DeleteProductCache
}

func NewDeleteProductBiz(store DeleteProductStorage, cache DeleteProductCache) *deleteProductBusiness {
	return &deleteProductBusiness{store, cache}
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

	var pagging paggingcommon.Paging
	pagging.Process()

	key := cacheutil.GenerateKey(productmodel.EntityName, cond, pagging, commonfilter.Filter{})

	if err := biz.cache.DeleteProduct(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	return nil
}
