package productstorage

import (
	"context"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (s *mysqlProduct) GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error) {

	db := s.db

	var record productmodel.Product

	if err := db.WithContext(ctx).Where(cond).
		Preload("Images").
		Preload("Category").
		Preload("Recipes").
		Preload("Recipes.RecipeIngredients").
		First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
