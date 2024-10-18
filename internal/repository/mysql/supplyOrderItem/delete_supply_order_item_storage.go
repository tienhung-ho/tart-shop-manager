package supplyorderitemstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrderItem) DeleteSupplyOrderItems(ctx context.Context, supplyOrderItemIDs []uint64) error {
	db := s.getDB(ctx)

	if err := db.WithContext(ctx).
		Where("supplyorderitem_id IN ?", supplyOrderItemIDs).
		Delete(&supplyordermodel.SupplyOrderItem{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Xử lý khi không tìm thấy bản ghi
			return err
		}
		return common.ErrDB(err)
	}

	return nil
}
