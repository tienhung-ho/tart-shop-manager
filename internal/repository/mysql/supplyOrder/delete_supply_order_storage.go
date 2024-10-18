package supplyorderstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrder) DeleteSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.getDB(ctx)

	if err := db.WithContext(ctx).
		Where(cond).
		Delete(&supplyordermodel.SupplyOrder{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Xử lý khi không tìm thấy bản ghi
			return err
		}
		return common.ErrDB(err)
	}

	return nil

	return nil
}
