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

//// GetTotalIngredients lấy tổng số lượng nguyên liệu cần thiết cho danh sách sản phẩm
//func (r *mysqlIngredient) GetTotalIngredients(ctx context.Context, products []ordermodel.Product) ([]ordermodel.RecipeIngredientQuantity, error) {
//	var results []ordermodel.RecipeIngredientQuantity
//
//	// Map để tổng hợp số lượng nguyên liệu
//	ingredientMap := make(map[int64]int64)
//
//	for _, p := range products {
//		var ingredients []ordermodel.RecipeIngredientQuantity
//		err := r.db.WithContext(ctx).Table("RecipeIngredients as ri").
//			Select("ri.ingredient_id, ri.quantity").
//			Joins("JOIN Recipe r ON ri.recipe_id = r.recipe_id").
//			Where("r.product_id = ? AND r.size = ?", p.ProductID, p.Size).
//			Scan(&ingredients).Error
//		if err != nil {
//			return nil, fmt.Errorf("failed to get ingredients for product %d and size %s: %v", p.ProductID, p.Size, err)
//		}
//
//		// Cộng số lượng nguyên liệu dựa trên số lượng sản phẩm
//		for _, ing := range ingredients {
//			ingredientMap[ing.IngredientID] += ing.Quantity * int64(p.Quantity)
//		}
//	}
//
//	// Chuyển đổi map thành slice
//	for id, qty := range ingredientMap {
//		results = append(results, ordermodel.RecipeIngredientQuantity{
//			IngredientID: id,
//			Quantity:     qty,
//		})
//	}
//
//	return results, nil
//}
