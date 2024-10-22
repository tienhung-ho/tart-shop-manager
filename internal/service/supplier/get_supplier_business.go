package supplierbusiness

import (
	"context"
	"errors"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetSupplierStorage interface {
	GetSupplier(ctx context.Context, cond map[string]interface{},
		morekeys ...string) (*suppliermodel.Supplier, error)
}

type GetSupplierCache interface {
	GetSupplier(ctx context.Context,
		cond map[string]interface{}, morekeys ...string) (*suppliermodel.Supplier, error)
	SaveSupplier(ctx context.Context, data interface{}, morekeys ...string) error
}

type getSupplierBusiness struct {
	store GetSupplierStorage
	cache GetSupplierCache
}

func NewGetSupplierBiz(store GetSupplierStorage, cache GetSupplierCache) *getSupplierBusiness {
	return &getSupplierBusiness{store, cache}
}

func (biz *getSupplierBusiness) GetSupplier(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*suppliermodel.Supplier, error) {

	record, err := biz.cache.GetSupplier(ctx, cond, morekeys...)
	if err != nil {
		return nil, common.ErrCannotGetEntity(suppliermodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetSupplier(ctx, cond, morekeys...)
	if err != nil {
		if errors.As(err, &common.RecordNotFound) {
			return nil, common.ErrNotFoundEntity(suppliermodel.EntityName, err)
		}

		return nil, common.ErrCannotGetEntity(suppliermodel.EntityName, err)
	}

	if record != nil {
		var paging paggingcommon.Paging
		paging.Process()

		createRecord := record.ToCreateSupplierCache()

		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: suppliermodel.EntityName,
			Cond:       cond,
			Paging:     paging,
			Filter:     commonfilter.Filter{},
			MoreKeys:   morekeys,
		})

		if err != nil {
			return nil, common.ErrCannotGenerateKey(suppliermodel.EntityName, err)
		}

		if err := biz.cache.SaveSupplier(ctx, createRecord, key); err != nil {
			return nil, common.ErrCannotCreateEntity(suppliermodel.EntityName, err)
		}

	}

	return record, nil
}
