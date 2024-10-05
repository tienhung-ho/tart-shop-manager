package recipestorage

import (
	"context"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

func (s *mysqlRecipe) GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error) {

	db := s.db

	var recipe recipemodel.Recipe
	if err := db.WithContext(ctx).Select(recipemodel.SelectFields).Where(cond).First(&recipe).Error; err != nil {

		return nil, err
	}

	return &recipe, nil
}
