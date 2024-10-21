package stockbatchbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteStockBatchStorage interface {
	GetStockBatch(ctx context.Context, cond map[string]interface{}) (*stockbatchmodel.StockBatch, error)
	DeleteStockBatch(ctx context.Context, cond map[string]interface{}) error
}

type DeleteStockBatchCache interface {
	DeleteStockBatch(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type deleteStockBatchCacheBusiness struct {
	store DeleteStockBatchStorage
	cache DeleteStockBatchCache
}

func NewDeleteStockBatchBiz(store DeleteStockBatchStorage,
	cache DeleteStockBatchCache) *deleteStockBatchCacheBusiness {
	return &deleteStockBatchCacheBusiness{store, cache}
}

func (biz *deleteStockBatchCacheBusiness) DeleteStockBatch(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) error {
	record, err := biz.store.GetStockBatch(ctx, cond)

	if err != nil {
		return common.ErrCannotGetEntity(stockbatchmodel.EntityName, err)
	}

	if record == nil {
		return common.ErrNotFoundEntity(stockbatchmodel.EntityName, err)
	}

	if err := biz.store.DeleteStockBatch(ctx, map[string]interface{}{"stockbatch_id": record.StockBatchID}); err != nil {
		return common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	var paging paggingcommon.Paging
	paging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: stockbatchmodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})

	if err != nil {
		return common.ErrInvalidGender(stockbatchmodel.EntityName, err)
	}

	if err := biz.cache.DeleteStockBatch(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, stockbatchmodel.EntityName); err != nil {
		return common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	return nil
}
