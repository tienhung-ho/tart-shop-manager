package supplierbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateSupplierStorage interface {
	UpdateSupplier(ctx context.Context, cond map[string]interface{},
		data *suppliermodel.UpdateSupplier) (*suppliermodel.Supplier, error)
	GetSupplier(ctx context.Context, cond map[string]interface{},
		morekeys ...string) (*suppliermodel.Supplier, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UpdateSupplierCache interface {
	DeleteSupplier(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type updateSupplierBusiness struct {
	store UpdateSupplierStorage
	cache UpdateSupplierCache
}

func NewUpdateSupplierBiz(store UpdateSupplierStorage, cache UpdateSupplierCache) *updateSupplierBusiness {
	return &updateSupplierBusiness{store, cache}
}

func (biz *updateSupplierBusiness) UpdateSupplier(ctx context.Context, cond map[string]interface{},
	data *suppliermodel.UpdateSupplier, morekeys ...string) (*suppliermodel.Supplier, error) {

	record, err := biz.store.GetSupplier(ctx, cond)

	if err != nil {
		return nil, common.ErrCannotGetEntity(suppliermodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(suppliermodel.EntityName, nil)
	}

	var updatedRecord *suppliermodel.Supplier

	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		updatedRecord, err = biz.store.UpdateSupplier(txCtx, cond, data)

		if err != nil {
			return common.ErrCannotUpdateEntity(suppliermodel.EntityName, err)
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

		if err := biz.cache.DeleteSupplier(txCtx, key); err != nil {
			return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
		}

		if err := biz.cache.DeleteListCache(ctx, suppliermodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
		}

		return nil
	})

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(suppliermodel.EntityName, err)
	}

	return updatedRecord, nil
}
