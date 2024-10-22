package supplierstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/common/trainsaction"
)

var (
	SelectFields = []string{
		"supplier_id",
		"name",
		"description",
		"address",
		"contactInfo",
		"status",
	}

	AllowedSortFields = map[string]bool{
		"created_at":  true,
		"updated_at":  true,
		"supplier_id": true,
	}
)

type mysqlSupplier struct {
	db *gorm.DB
}

func NewMySQLSupplier(db *gorm.DB) *mysqlSupplier {
	return &mysqlSupplier{db}
}

func (s *mysqlSupplier) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
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
func (s *mysqlSupplier) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return s.db
}
