package productbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
}

type GetProductCache interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	SaveProduct(ctx context.Context, data interface{}, morekeys ...string) error
}

type getProductBusiness struct {
	store GetProductStorage
	cache GetProductCache
}

func NewGetProductBiz(store GetProductStorage, cache GetProductCache) *getProductBusiness {
	return &getProductBusiness{store: store, cache: cache}
}

func (biz *getProductBusiness) GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error) {

	record, err := biz.cache.GetProduct(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetProduct(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
	}

	if record != nil {
		var pagging paggingcommon.Paging
		pagging.Process()

		var createProduct = record.ToCreateAccount()

		key := cacheutil.GenerateKey(productmodel.EntityName, cond, pagging, commonfilter.Filter{})

		if err := biz.cache.SaveProduct(ctx, createProduct, key); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}

	}

	return record, nil
}
