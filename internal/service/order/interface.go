package orderbusiness

import (
	"context"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

type ProductStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error)
}

type IngredientStorage interface {
	GetIngredientIDsByProductID(ctx context.Context,
		productIDs []uint64) ([]ordermodel.RecipeIngredientQuantity, error)
}

type RecipeStorage interface {
	GetRecipe(ctx context.Context,
		cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]recipemodel.Recipe, error)
}

type StockBatchStorage interface {
	GetStockBatchesByIngredientIDs(ctx context.Context,
		ingredientIDs []uint64) ([]stockbatchmodel.StockBatch, error)
	UpdateStockBatches(ctx context.Context, cond map[string]interface{},
		data []stockbatchmodel.UpdateStockBatch) ([]uint64, error)
}
