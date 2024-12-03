package supplyorderitemstorage

import (
	"context"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrderItem) ListItemByIngredients(ctx context.Context,
	ingredientIDs []uint64) ([]supplyordermodel.SimpleSupplyOrderItemIngredient, error) {
	var supplyOrderItems []supplyordermodel.SimpleSupplyOrderItemIngredient

	// Thực hiện tìm kiếm với danh sách ingredientIDs
	if err := s.db.WithContext(ctx).
		Table(supplyordermodel.SimpleSupplyOrderItemIngredient{}.TableName()).
		Where("ingredient_id IN ?", ingredientIDs).
		Find(&supplyOrderItems).Error; err != nil {
		return nil, err
	}

	return supplyOrderItems, nil
}
