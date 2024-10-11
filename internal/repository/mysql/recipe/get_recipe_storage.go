package recipestorage

import (
	"context"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

func (s *mysqlRecipe) GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error) {

	db := s.db

	var recipe recipemodel.Recipe
	if err := db.WithContext(ctx).Select(SelectFields).Where(cond).Preload("RecipeIngredients").Preload("Product").First(&recipe).Error; err != nil {

		return nil, err
	}

	return &recipe, nil
}
