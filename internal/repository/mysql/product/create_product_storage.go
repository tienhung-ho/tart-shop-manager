package productstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlProduct) CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.ProductID, nil
}
