package categorystorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlCategory) CreateCategory(ctx context.Context, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.CategoryID, nil
}
