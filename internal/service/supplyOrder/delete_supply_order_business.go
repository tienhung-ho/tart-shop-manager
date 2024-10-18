package supplyorderbusiness

import (
	"context"
	"errors"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteSupplyOrderStorage interface {
	DeleteSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
	GetSupplyOrder(ctx context.Context, cond map[string]interface{}) (*supplyordermodel.SupplyOrder, error)
}

type DeleteSupplyOrderCache interface {
	DeleteSupplyOrder(ctx context.Context, morekeys ...string) error
}

type deleteSupplyOrderBusiness struct {
	store           DeleteSupplyOrderStorage
	cache           DeleteSupplyOrderCache
	storeItem       CreateSupplyOrderItemStorage
	storeIngredient IngredientStorage
	storeStockBatch StockBatchStorage
}

func NewDeleteSupplyOrderBiz(store DeleteSupplyOrderStorage,
	cache DeleteSupplyOrderCache,
	storeItem CreateSupplyOrderItemStorage,
	storeIngredient IngredientStorage,
	storeStockBatch StockBatchStorage) *deleteSupplyOrderBusiness {
	return &deleteSupplyOrderBusiness{store, cache, storeItem, storeIngredient, storeStockBatch}
}

func (biz *deleteSupplyOrderBusiness) DeleteSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetSupplyOrder(ctx, cond)

	if err != nil {
		if errors.Is(err, common.RecordNotFound) {
			return common.ErrNotFoundEntity(supplyordermodel.EntityName, err)
		}
		return common.ErrCannotGetEntity(supplyordermodel.EntityName, err)
	}

	if record == nil {
		return common.ErrNotFoundEntity(supplyordermodel.EntityName, err)
	}

	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		if err := biz.store.DeleteSupplyOrder(txCtx, cond, morekeys...); err != nil {
			return common.ErrCannotDeleteEntity(supplyordermodel.EntityName, err)
		}

		var itemsToDeleteIDs []uint64
		var stockBatchDeleteIDs []uint64
		for _, item := range record.SupplyOrderItems {
			itemsToDeleteIDs = append(itemsToDeleteIDs, item.SupplyOrderItemID)
			stockBatchDeleteIDs = append(stockBatchDeleteIDs, item.StockBatchID)
		}

		if len(itemsToDeleteIDs) > 0 {
			if err := biz.storeItem.DeleteSupplyOrderItems(txCtx, itemsToDeleteIDs); err != nil {
				return common.ErrCannotDeleteEntity(supplyordermodel.ItemEntityName, err)
			}
		}

		if len(stockBatchDeleteIDs) > 0 {
			if err := biz.storeStockBatch.DeleteStockBatches(txCtx, stockBatchDeleteIDs); err != nil {
				return common.ErrCannotUpdateEntity(stockbatchmodel.EntityName, err)
			}
		}

		var pagging paggingcommon.Paging
		pagging.Process()

		// Generate cache key
		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: supplyordermodel.EntityName,
			Cond:       cond,
			Paging:     pagging,
			Filter:     commonfilter.Filter{},
			MoreKeys:   morekeys,
		})
		if err != nil {
			return common.ErrCannotGenerateKey(supplyordermodel.EntityName, err)
		}

		if err := biz.cache.DeleteSupplyOrder(ctx, key); err != nil {
			return common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
		}

		return nil
	})

	return nil
}
