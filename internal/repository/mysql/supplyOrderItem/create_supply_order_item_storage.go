package supplyorderitemstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrderItem) CreateSupplyOrderItem(ctx context.Context, item []supplyordermodel.CreateSupplyOrderItem) error {
	if len(item) == 0 {
		return nil
	}

	// Sử dụng Bulk Insert để tối ưu hóa hiệu suất
	if err := s.getDB(ctx).WithContext(ctx).Create(&item).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
