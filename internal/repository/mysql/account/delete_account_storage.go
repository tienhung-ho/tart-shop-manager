package accountstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
)

func (s *mysqlAccount) DeleteAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {
	db := s.db

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	var account accountmodel.Account
	if err := db.WithContext(ctx).Where(cond).Delete(&account).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
