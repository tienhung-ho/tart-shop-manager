package supplierbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

type CreateSupplierStorage interface {
	CreateSupplier(ctx context.Context, data *suppliermodel.CreateSupplier) (uint64, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type CreateSupplierCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type createSupplierBusiness struct {
	store CreateSupplierStorage
	cache CreateSupplierCache
}

func NewCreateSupplierBusiness(store CreateSupplierStorage, cache CreateSupplierCache) *createSupplierBusiness {
	return &createSupplierBusiness{store: store, cache: cache}
}

func (biz *createSupplierBusiness) CreateSupplier(ctx context.Context, data *suppliermodel.CreateSupplier) (uint64, error) {

	var recordID uint64
	var err error
	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		recordID, err = biz.store.CreateSupplier(txCtx, data)

		if err != nil {
			return common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
		}

		if err := biz.cache.DeleteListCache(ctx, suppliermodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(suppliermodel.EntityName, err)
		}

		return nil
	})

	return recordID, err
}
