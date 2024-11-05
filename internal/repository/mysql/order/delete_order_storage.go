package orderstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

func (s *mysqlOrder) DeleteOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	tx := s.getDB(ctx)

	if err := tx.Where(cond).Delete(&ordermodel.Order{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Xử lý khi không tìm thấy bản ghi
			return common.RecordNotFound
		}
		return common.ErrDB(err)
	}

	return nil
}
