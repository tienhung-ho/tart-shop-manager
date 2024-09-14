package accountstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	"tart-shop-manager/internal/entity/dtos/sql/account"
)

func (s *mysqlAccount) CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error) {

	db := s.db.Begin()

	defer commonrecover.RecoverTransaction(db)

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		db.Rollback()
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		return 0, common.ErrDB(err)
	}

	return data.AccountID, nil
}
