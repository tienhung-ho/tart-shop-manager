package orderitemstorage

import (
	"context"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

func (s *mysqlOrderItem) DeleteOrderItem(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	tx := s.getDB(ctx)

	if err := tx.Where(cond).Delete(&ordermodel.OrderItem{}).Error; err != nil {

		return err
	}

	return nil
}
