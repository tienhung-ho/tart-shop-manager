package supplierstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
)

func (s *mysqlSupplier) GetSupplier(ctx context.Context, cond map[string]interface{},
	morekeys ...string) (*suppliermodel.Supplier, error) {

	db := s.db

	var record suppliermodel.Supplier
	if err := db.WithContext(ctx).Where(cond).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
