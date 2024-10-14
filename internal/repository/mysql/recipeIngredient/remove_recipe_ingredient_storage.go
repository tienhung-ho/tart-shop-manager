package recipeingredientstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

func (s *mysqlRecipeIngredient) RemoveRecipeIngredients(ctx context.Context, recipeID uint64, ingredientIDs []uint64) error {
	if err := s.getDB(ctx).
		Where("recipe_id = ? AND ingredient_id IN ?", recipeID, ingredientIDs).
		Delete(&recipeingredientmodel.RecipeIngredient{}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
