package recipebusiness

import (
	"context"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

type ListItemIngredientStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, moreKeys ...string) ([]ingredientmodel.Ingredient, error)
}

type RecipeIngredientStorage interface {
	CreateRecipeIngredients(ctx context.Context, ingredients []recipeingredientmodel.RecipeIngredientCreate) error
	RemoveRecipeIngredients(ctx context.Context, recipeID uint64, ingredientIDs []uint64) error
	UpdateRecipeIngredient(ctx context.Context, ing *recipeingredientmodel.UpdateRecipeIngredient) error
}
