package supplierbusiness

import (
	"context"
	"fmt"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemSupplierStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]suppliermodel.Supplier, error)
}

type ListItemSupplierCache interface {
	ListItem(ctx context.Context, key string) ([]suppliermodel.Supplier, error)
	SaveSupplier(ctx context.Context, data interface{}, morekeys ...string) error
}

type listItemSupplierBusiness struct {
	store ListItemSupplierStorage
	cache ListItemSupplierCache
}

func NewListItemSupplierBiz(store ListItemSupplierStorage, cache ListItemSupplierCache) *listItemSupplierBusiness {
	return &listItemSupplierBusiness{store, cache}
}

func (biz *listItemSupplierBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]suppliermodel.Supplier, error) {

	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: suppliermodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", suppliermodel.EntityName),
	})

	if err != nil {
		return nil, common.ErrCannotGenerateKey(suppliermodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, key)

	if err != nil {
		return nil, common.ErrCannotListEntity(suppliermodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, &pagingCopy, &filterCopy, morekeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(suppliermodel.EntityName, err)
	}

	if len(records) != 0 {
		if err := biz.cache.SaveSupplier(ctx, records, key); err != nil {
		}

		return records, nil
	}

	return nil, common.ErrCannotListEntity(suppliermodel.EntityName, err)
}
