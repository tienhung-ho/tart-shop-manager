package orderbusiness

import (
	"context"
	"fmt"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListOrderStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]ordermodel.Order, error)
}

type ListOrderCache interface {
	ListItem(ctx context.Context, key string) ([]ordermodel.Order, error)
	SaveOrder(ctx context.Context, data interface{}, morekeys ...string) error
}

type listOrderBusiness struct {
	store ListOrderStorage
	cache ListOrderCache
}

func NewListOrderBiz(store ListOrderStorage, cache ListOrderCache) *listOrderBusiness {
	return &listOrderBusiness{store, cache}
}

func (biz *listOrderBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]ordermodel.Order, error) {

	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: ordermodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", ordermodel.EntityName),
	})

	if err != nil {
		return nil, common.ErrCannotGenerateKey(ordermodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, key)

	if err != nil {
		return nil, common.ErrCannotListEntity(ordermodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, paging, filter, morekeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(ordermodel.EntityName, err)
	}

	if records == nil {
		return nil, common.ErrCannotListEntity(ordermodel.EntityName, err)
	}

	if err := biz.cache.SaveOrder(ctx, records, key); err != nil {
		return nil, common.ErrCannotCreateEntity(ordermodel.EntityName, err)
	}

	return records, nil
}
