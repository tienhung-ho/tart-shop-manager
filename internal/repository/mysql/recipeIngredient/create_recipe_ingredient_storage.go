package recipeingredientstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

// CreateRecipeIngredients thêm danh sách nguyên liệu vào công thức
func (s *mysqlRecipeIngredient) CreateRecipeIngredients(ctx context.Context, ingredients []recipeingredientmodel.RecipeIngredientCreate) error {
	if len(ingredients) == 0 {
		return nil
	}

	// Sử dụng Bulk Insert để tối ưu hóa hiệu suất
	if err := s.getDB(ctx).WithContext(ctx).Create(&ingredients).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
