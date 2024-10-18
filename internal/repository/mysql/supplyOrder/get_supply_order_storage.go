package supplyorderstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrder) GetSupplyOrder(ctx context.Context, cond map[string]interface{}) (*supplyordermodel.SupplyOrder, error) {

	db := s.db

	var record supplyordermodel.SupplyOrder
	if err := db.WithContext(ctx).Select(SelectFields).Where(cond).
		Preload("SupplyOrderItems").
		Preload("SupplyOrderItems.StockBatch").
		First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Xử lý khi không tìm thấy bản ghi
			return nil, err
		}
		return nil, err
	}

	return &record, nil
}
