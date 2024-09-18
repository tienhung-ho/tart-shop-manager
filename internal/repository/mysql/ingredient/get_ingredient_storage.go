package ingredientstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error) {

	db := s.db

	var result ingredientmodel.Ingredient
	if err := db.WithContext(ctx).Where(cond).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound // Trả về nil nếu không tìm thấy
		}
		return nil, err
	}

	return &result, nil
}
