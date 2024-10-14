package recipeingredientstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common/trainsaction"
)

var ()

type mysqlRecipeIngredient struct {
	db *gorm.DB
}

func NewMySQLRecipeIngredient(db *gorm.DB) *mysqlRecipeIngredient {
	return &mysqlRecipeIngredient{db}
}

// Hàm hỗ trợ lấy *gorm.DB từ context nếu có transaction
func (s *mysqlRecipeIngredient) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(trainsaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return s.db
}
