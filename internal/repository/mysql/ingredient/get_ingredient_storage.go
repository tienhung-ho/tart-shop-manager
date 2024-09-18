package ingredientstorage

import (
	"context"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error) {

	db := s.db

	var result ingredientmodel.Ingredient
	if err := db.WithContext(ctx).Where(cond).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
