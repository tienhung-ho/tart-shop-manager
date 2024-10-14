package recipeingredientstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

func (s *mysqlRecipeIngredient) UpdateRecipeIngredient(ctx context.Context, ing *recipeingredientmodel.UpdateRecipeIngredient) error {
	if err := s.getDB(ctx).
		Model(&recipeingredientmodel.UpdateRecipeIngredient{}).
		Where("recipe_id = ? AND ingredient_id = ?", ing.RecipeID, ing.IngredientID).
		Update("quantity", ing.Quantity).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
