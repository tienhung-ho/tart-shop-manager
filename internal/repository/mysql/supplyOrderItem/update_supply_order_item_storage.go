package supplyorderitemstorage

import (
	"context"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrderItem) UpdateSupplyOrderItems(ctx context.Context, items []supplyordermodel.UpdateSupplyOrderItem) error {

	for _, item := range items {
		err := s.db.WithContext(ctx).Model(&supplyordermodel.SupplyOrderItem{}).
			Where("supplyorderitem_id = ?", item.SupplyOrderItemID).
			Updates(map[string]interface{}{
				"stockbatch_id": item.StockBatchID,
				"price":         item.Price,
				"quantity":      item.Quantity,
				"unit":          item.Unit,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
