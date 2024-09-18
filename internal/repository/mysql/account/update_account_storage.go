package accountstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlAccount) UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&accountmodel.CreateAccount{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(account).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		}

		db.Rollback()
		return nil, err
	}

	// Truy vấn lại bản ghi đã cập nhật
	var updatedRecord accountmodel.Account
	if err := db.WithContext(ctx).Model(&accountmodel.Account{}).Where(cond).Preload("Role").First(&updatedRecord).Error; err != nil {
		db.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &updatedRecord, nil
}
