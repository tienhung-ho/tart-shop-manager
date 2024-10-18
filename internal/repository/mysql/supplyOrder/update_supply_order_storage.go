package supplyorderstorage

import (
	"context"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (s *mysqlSupplyOrder) UpdateSupplyOrder(ctx context.Context, cond map[string]interface{}, data *supplyordermodel.UpdateSupplyOrder) (*supplyordermodel.SupplyOrder, error) {

	return nil, nil
}
