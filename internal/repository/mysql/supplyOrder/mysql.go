package supplyorderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/common/trainsaction"
)

var (
	SelectFields = []string{
		"supplyorder_id",
		"order_date",
		"description",
		"total_amount",
		"supplier_id",
		"status",
	}

	AllowedSortFields = map[string]bool{
		"order_date":     true,
		"created_at":     true,
		"updated_at":     true,
		"supplyorder_id": true,
	}
)

type mysqlSupplyOrder struct {
	db *gorm.DB
}

func NewMySQLSupplyOrder(db *gorm.DB) *mysqlSupplyOrder {
	return &mysqlSupplyOrder{db}
}

func (r *mysqlSupplyOrder) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
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
func (r *mysqlSupplyOrder) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}
