package stockbatchbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetStockBatchStorage interface {
	GetStockBatch(ctx context.Context, cond map[string]interface{}) (*stockbatchmodel.StockBatch, error)
}

type GetStockBatchCache interface {
	GetStockBatch(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*stockbatchmodel.StockBatch, error)
	SaveStockBatch(ctx context.Context, data interface{}, morekeys ...string) error
}

type getStockBatchBusiness struct {
	store GetStockBatchStorage
	cache GetStockBatchCache
}

func NewGetStockBatchBiz(store GetStockBatchStorage, cache GetStockBatchCache) *getStockBatchBusiness {
	return &getStockBatchBusiness{store, cache}
}

func (biz *getStockBatchBusiness) GetStockBatch(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*stockbatchmodel.StockBatch, error) {

	record, err := biz.cache.GetStockBatch(ctx, cond)

	if err != nil {
		return nil, common.ErrCannotGetEntity(stockbatchmodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetStockBatch(ctx, cond)

	if err != nil {
		return nil, common.ErrCannotGetEntity(stockbatchmodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(stockbatchmodel.EntityName, nil)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	createStockBatch := record.ToCreateStockBatch()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: stockbatchmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	
	if err != nil {
		return nil, common.ErrCannotGenerateKey(supplyordermodel.EntityName, err)
	}

	if err := biz.cache.SaveStockBatch(ctx, createStockBatch, key); err != nil {
		return nil, common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
	}

	return record, nil
}
