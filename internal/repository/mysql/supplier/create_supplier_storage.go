package supplierstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlSupplier) CreateSupplier(ctx context.Context, data *suppliermodel.CreateSupplier) (uint64, error) {

	db := s.getDB(ctx)

	if err := db.Create(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := responseutil.ExtractFieldFromError(err, suppliermodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(suppliermodel.EntityName, fieldName, err)
		}
		return 0, err
	}

	return data.SupplierID, nil
}
