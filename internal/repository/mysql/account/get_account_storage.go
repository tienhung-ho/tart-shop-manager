package accountstorage

import (
	"context"
	"tart-shop-manager/internal/entity/dtos/sql/account"
)

func (s *mysqlAccount) GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error) {

	db := s.db

	//defer commonrecover.RecoverTransaction(db)
	var record accountmodel.Account

	if err := db.WithContext(ctx).
		Select(accountmodel.SelectFields).
		Where(cond).Preload("Images").Preload("Role").First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
