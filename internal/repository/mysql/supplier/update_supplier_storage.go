package supplierstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
)

func (s *mysqlSupplier) UpdateSupplier(ctx context.Context, cond map[string]interface{},
	data *suppliermodel.UpdateSupplier) (*suppliermodel.Supplier, error) {

	db := s.getDB(ctx)

	if err := db.WithContext(ctx).Where(cond).Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}

		return nil, err
	}

	var record suppliermodel.Supplier
	if err := db.Where(cond).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
