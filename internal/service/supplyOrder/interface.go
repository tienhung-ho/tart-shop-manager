package supplyorderbusiness

import (
	"context"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	"tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

type IngredientStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paggingcommon.Paging, filter *commonfilter.Filter,
		moreKeys ...string) ([]ingredientmodel.Ingredient, error)
}

type CreateSupplyOrderItemStorage interface {
	CreateSupplyOrderItems(ctx context.Context, item []supplyordermodel.CreateSupplyOrderItem) error
	UpdateSupplyOrderItems(ctx context.Context, items []supplyordermodel.UpdateSupplyOrderItem) error
	DeleteSupplyOrderItems(ctx context.Context, supplyOrderItemIDs []uint64) error
}

type StockBatchStorage interface {
	CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint64, error)
	CreateStockBatches(ctx context.Context, data []stockbatchmodel.CreateStockBatch) ([]uint64, error)
	UpdateStockBatches(ctx context.Context, cond map[string]interface{}, data []stockbatchmodel.UpdateStockBatch) ([]uint64, error)
	GetStockBatch(ctx context.Context, cond map[string]interface{}) (*stockbatchmodel.StockBatch, error)
	DeleteStockBatches(ctx context.Context, stockBatchIDs []uint64) error
}
