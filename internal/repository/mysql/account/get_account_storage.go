package accountstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	"tart-shop-manager/internal/entity/model/sql/account"
)

func (s *mysqlAccount) GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error) {

	db := s.db

	defer commonrecover.RecoverTransaction(db)
	var record accountmodel.Account

	if err := db.WithContext(ctx).Where(cond).First(&record).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
