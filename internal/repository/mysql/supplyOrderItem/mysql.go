package supplyorderitemstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/common/trainsaction"
)

type mysqlSupplyOrderItem struct {
	db *gorm.DB
}

func NewMySQLSupplyOrderItem(db *gorm.DB) *mysqlSupplyOrderItem {
	return &mysqlSupplyOrderItem{db}
}

func (s *mysqlSupplyOrderItem) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return common.ErrDB(tx.Error)
	}

	txCtx := context.WithValue(ctx, trainsaction.TransactionKey, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return common.ErrDB(err)
	}

	return nil
}

// getDB lấy *gorm.DB từ context nếu có transaction
func (s *mysqlSupplyOrderItem) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return s.db
}
