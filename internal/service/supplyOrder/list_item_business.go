package supplyorderbusiness

import (
	"context"
	"fmt"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemSupplyOrderStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error)
}

type ListItemSupplyOrderCache interface {
	ListItem(ctx context.Context, key string) ([]supplyordermodel.SupplyOrder, error)
	DeleteListCache(ctx context.Context, entityName string) error
	SaveSupplyOrder(ctx context.Context, data interface{}, morekeys ...string) error
}

type listItemSupplyOrderBusiness struct {
	store ListItemSupplyOrderStorage
	cache ListItemSupplyOrderCache
}

func NewListItemSupplyOrderBiz(store ListItemSupplyOrderStorage, cache ListItemSupplyOrderCache) *listItemSupplyOrderBusiness {
	return &listItemSupplyOrderBusiness{store: store, cache: cache}
}

func (biz *listItemSupplyOrderBusiness) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error) {

	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: supplyordermodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", supplyordermodel.EntityName),
	})

	records, err := biz.cache.ListItem(ctx, key)

	if err != nil {
		return nil, common.ErrCannotListEntity(supplyordermodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, &pagingCopy, &filterCopy, morekeys...)
	if err != nil {
		return nil, common.ErrCannotListEntity(supplyordermodel.EntityName, err)
	}

	if len(records) != 0 {
		if err := biz.cache.SaveSupplyOrder(ctx, records, key); err != nil {
			log.Print(err)
			return nil, common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
		}
	}

	return records, nil

}
