package recipestorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/common/trainsaction"
)

var (
	SelectFields = []string{
		"recipe_id",
		"product_id",
		"size",
		"cost",
		"description",
		"status",
	}
	AllowedSortFields = map[string]bool{
		"recipe_id":  true,
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"cost":       true,
	}
)

type mysqlRecipe struct {
	db *gorm.DB
}

func NewMySQLRecipe(db *gorm.DB) *mysqlRecipe {
	return &mysqlRecipe{db}
}

func (s *mysqlRecipe) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
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
func (s *mysqlRecipe) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return s.db
}
