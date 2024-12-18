package orderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/common/trainsaction"
)

var (
	AllowedSortFields = map[string]bool{
		"order_id":   true,
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"cost":       true,
	}
)

type mysqlOrder struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlOrder {
	return &mysqlOrder{db}
}

func (r *mysqlOrder) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := r.db.Begin()
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
func (r *mysqlOrder) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}
