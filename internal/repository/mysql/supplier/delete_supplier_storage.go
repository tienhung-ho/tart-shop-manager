package supplierstorage

import (
	"context"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
)

func (s *mysqlSupplier) DeleteSupplier(ctx context.Context, cond map[string]interface{}) error {

	db := s.db

	if err := db.WithContext(ctx).Where(cond).Delete(&suppliermodel.Supplier{}).Error; err != nil {
		return err
	}

	return nil
}
