package supplyorderbusiness

import (
	"context"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

type IngredientStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paggingcommon.Paging, filter *commonfilter.Filter,
		moreKeys ...string) ([]ingredientmodel.Ingredient, error)
}

type CreateSupplyOrderItemStorage interface {
	CreateSupplyOrderItem(ctx context.Context, data []supplyordermodel.CreateSupplyOrderItem) error
}

type StockBatchStorage interface {
	CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint, error)
}
