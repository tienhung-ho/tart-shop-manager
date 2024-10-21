package stockbatchbusiness

import (
	"context"
	"fmt"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]stockbatchmodel.StockBatch, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, key string) ([]stockbatchmodel.StockBatch, error)
	SaveStockBatch(ctx context.Context, data interface{}, morekeys ...string) error
}

type listItemStockBatchBusiness struct {
	store ListItemStorage
	cache ListItemCache
}

func NewListItemStockBatchBiz(store ListItemStorage, cache ListItemCache) *listItemStockBatchBusiness {
	return &listItemStockBatchBusiness{store, cache}
}
func (biz *listItemStockBatchBusiness) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]stockbatchmodel.StockBatch, error) {

	pagingCopy := *pagging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: stockbatchmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", stockbatchmodel.EntityName),
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(stockbatchmodel.EntityName, err)
	}

	// Gọi cache với khóa đã tạo
	records, err := biz.cache.ListItem(ctx, key)
	if err != nil {
		return nil, common.ErrCannotListEntity(stockbatchmodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	// Gọi store để lấy dữ liệu từ database
	records, err = biz.store.ListItem(ctx, cond, pagging, filter, morekeys...)
	if err != nil {
		log.Print(err)
		return nil, common.ErrCannotListEntity(stockbatchmodel.EntityName, err)
	}

	// Lưu vào cache với cùng khóa
	if len(records) != 0 {
		if err := biz.cache.SaveStockBatch(ctx, records, key); err != nil {
			return nil, common.ErrCannotCreateEntity(stockbatchmodel.EntityName, err)
		}
	}

	return records, nil
}
