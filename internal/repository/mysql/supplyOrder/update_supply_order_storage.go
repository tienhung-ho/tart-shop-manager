package supplyorderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

// UpdateSupplyOrder updates the SupplyOrder record based on conditions and data provided
func (s *mysqlSupplyOrder) UpdateSupplyOrder(ctx context.Context, cond map[string]interface{}, data *supplyordermodel.UpdateSupplyOrder) (*supplyordermodel.SupplyOrder, error) {
	var supplyOrder supplyordermodel.SupplyOrder
	if err := s.db.WithContext(ctx).Where(cond).First(&supplyOrder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrNotFoundEntity(supplyordermodel.EntityName, nil)
		}
		return nil, common.ErrDB(err)
	}

	var updatedBy string
	email, ok := ctx.Value("email").(string)

	if !ok {
		updatedBy = "system" // Hoặc giá trị mặc định khác
	}

	updatedBy = email

	// Cập nhật các trường chính của SupplyOrder
	updateData := map[string]interface{}{
		"order_date":   data.OrderDate,
		"description":  data.Description,
		"total_amount": data.TotalAmount,
		"supplier_id":  data.SupplierID,
		"updated_by":   updatedBy,
		// Thêm các trường khác nếu có
	}

	if err := s.db.WithContext(ctx).Model(&supplyOrder).Updates(&updateData).Error; err != nil {
		return nil, common.ErrCannotUpdateEntity(supplyordermodel.EntityName, err)
	}

	// Lấy lại SupplyOrder đã cập nhật
	if err := s.db.WithContext(ctx).
		Preload("SupplyOrderItems").
		Preload("SupplyOrderItems.StockBatch").
		First(&supplyOrder, "supplyorder_id = ?", supplyOrder.SupplyOrderID).Error; err != nil {
		return nil, common.ErrCannotGetEntity(supplyordermodel.EntityName, err)
	}

	return &supplyOrder, nil
}
