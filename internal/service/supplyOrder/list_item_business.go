package supplyorderbusiness

import (
	"context"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

type GetItemSupplyOrderStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error)
}

type GetItemSupplyOrderCache interface {
	ListItem(ctx context.Context, key string) ([]supplyordermodel.SupplyOrder, error)
	DeleteListCache(ctx context.Context, entityName string) error
}

type getItemSupplyOrderBusiness struct {
	store GetItemSupplyOrderStorage
	cache GetItemSupplyOrderCache
}

func NewGetItemSupplyOrderBiz(store GetItemSupplyOrderStorage, cache GetItemSupplyOrderCache) *getItemSupplyOrderBusiness {
	return &getItemSupplyOrderBusiness{store: store, cache: cache}
}

func (biz *getItemSupplyOrderBusiness) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error) {

	//pagingCopy := *paging
	//filterCopy := *filter
	//
	//// Tạo khóa cache
	//key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
	//	EntityName: supplyordermodel.EntityName,
	//	Cond:       cond,
	//	Paging:     pagingCopy,
	//	Filter:     filterCopy,
	//	MoreKeys:   morekeys,
	//	KeyType:    fmt.Sprintf("List:%s:", supplyordermodel.EntityName),
	//})

	//records, err := biz.cache.ListItem(ctx, key)
	//
	//if err != nil {
	//	return nil, common.ErrCannotListEntity(supplyordermodel.EntityName, err)
	//}
	//
	//

	return nil, nil

}
