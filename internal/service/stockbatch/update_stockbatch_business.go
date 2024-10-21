package stockbatchbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateStockBatchStorage interface {
	GetStockBatch(ctx context.Context, cond map[string]interface{}) (*stockbatchmodel.StockBatch, error)
	UpdateStockBatch(ctx context.Context, cond map[string]interface{},
		data *stockbatchmodel.UpdateStockBatch) (*stockbatchmodel.StockBatch, error)
}

type UpdateStockBatchCache interface {
	DeleteStockBatch(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type updateStockBatchBusiness struct {
	store UpdateStockBatchStorage
	cache UpdateStockBatchCache
}

func NewUpdateStockBatchBiz(store UpdateStockBatchStorage, cache UpdateStockBatchCache) *updateStockBatchBusiness {
	return &updateStockBatchBusiness{store, cache}
}

func (biz *updateStockBatchBusiness) UpdateStockBatch(ctx context.Context, cond map[string]interface{},
	data *stockbatchmodel.UpdateStockBatch, morekeys ...string) (*stockbatchmodel.StockBatch, error) {

	record, err := biz.store.GetStockBatch(ctx, cond)

	if err != nil {
		return nil, common.ErrCannotGetEntity(stockbatchmodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(stockbatchmodel.EntityName, err)
	}

	update, err := biz.store.UpdateStockBatch(ctx, map[string]interface{}{"stockbatch_id": record.StockBatchID}, data)
	if err != nil {
		return nil, common.ErrCannotUpdateEntity(stockbatchmodel.EntityName, err)
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
		return nil, common.ErrCannotGenerateKey(stockbatchmodel.EntityName, err)
	}

	if err := biz.cache.DeleteStockBatch(ctx, key); err != nil {
		return nil, common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, stockbatchmodel.EntityName); err != nil {
		return nil, common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	return update, nil
}
