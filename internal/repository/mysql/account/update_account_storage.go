package accountstorage

import (
	"context"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
)

func (s *mysqlAccount) UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&accountmodel.Account{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(account).Error; err != nil {

		db.Rollback()
		return nil, err
	}

	// Truy vấn lại bản ghi đã cập nhật
	var updatedRecord accountmodel.Account
	if err := db.WithContext(ctx).Model(&accountmodel.Account{}).Where(cond).First(&updatedRecord).Error; err != nil {
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
