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

type DeleteSupplierStorage interface {
	GetSupplier(ctx context.Context, cond map[string]interface{},
		morekeys ...string) (*suppliermodel.Supplier, error)
	DeleteSupplier(ctx context.Context, cond map[string]interface{}) error
}

type DeleteSupplierCache interface {
	DeleteSupplier(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type deleteSupplierBusiness struct {
	store DeleteSupplierStorage
	cache DeleteSupplierCache
}

func NewDeleteSupplierBiz(s DeleteSupplierStorage, c DeleteSupplierCache) *deleteSupplierBusiness {
	return &deleteSupplierBusiness{store: s, cache: c}
}

func (biz *deleteSupplierBusiness) DeleteSupplier(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetSupplier(ctx, cond)

	if err != nil {
		if errors.Is(err, common.RecordNotFound) {
			return common.ErrNotFoundEntity(suppliermodel.EntityName, err)
		}
		return common.ErrCannotGetEntity(suppliermodel.EntityName, err)
	}

	if err := biz.store.DeleteSupplier(ctx, map[string]interface{}{"supplier_id": record.SupplierID}); err != nil {
		return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
	}

	var paging paggingcommon.Paging
	paging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: suppliermodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})

	if err != nil {
		return common.ErrCannotGenerateKey(suppliermodel.EntityName, err)
	}

	if err := biz.cache.DeleteSupplier(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, suppliermodel.EntityName); err != nil {
		return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
	}

	return nil
}
