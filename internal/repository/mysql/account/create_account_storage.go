package accountstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	"tart-shop-manager/internal/entity/dtos/sql/account"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlAccount) CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error) {

	db := s.db.Begin()

	defer commonrecover.RecoverTransaction(db)

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

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
		return 0, common.ErrDB(err)
	}

	return data.AccountID, nil
}
