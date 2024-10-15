package supplyorderstorage

import (
	"context"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrder) GetSupplyOrder(ctx context.Context, cond map[string]interface{}) (*supplyordermodel.SupplyOrder, error) {

	db := s.db

	var record supplyordermodel.SupplyOrder
	if err := db.WithContext(ctx).Select(SelectFields).Where(cond).
		Preload("SupplyOrderItems").
		Preload("SupplyOrderItems.StockBatch").
		First(&record).Error; err != nil {

	}

	return &record, nil
}

git commit -m "feat: Implement Get Supply Order with Redis Caching

- Add new handler for fetching supply orders
- Introduce Redis DTO and caching mechanisms
- Create MySQL storage for supply orders
- Update SupplyOrderItem and StockBatch models
- Add repository and service layers for Get Supply Order functionality
- Modify router and service interfaces for integration
- Enhance StockBatch storage with create functionality"
