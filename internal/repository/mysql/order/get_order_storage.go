package orderstorage

import (
	"context"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

func (s *mysqlOrder) GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error) {

	db := s.db

	var order ordermodel.Order
	if err := db.WithContext(ctx).Where(cond).
		Preload("OrderItems").
		Preload("OrderItems.Recipe").
		Preload("OrderItems.Recipe.Product").
		First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
