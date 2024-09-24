package ingredientstorage

import (
	"context"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) DeleteIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db
	var record ingredientmodel.Ingredient
	if err := db.WithContext(ctx).Where(cond).Delete(&record).Error; err != nil {
		return err
	}
	
	return nil
}
