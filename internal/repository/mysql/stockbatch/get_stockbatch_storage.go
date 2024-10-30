package stockbatchstorage

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (s *mysqlStockBatch) GetStockBatch(ctx context.Context, cond map[string]interface{}) (*stockbatchmodel.StockBatch, error) {
	var stockBatch stockbatchmodel.StockBatch

	if err := s.db.WithContext(ctx).
		Where(cond).
		Preload("Ingredient").
		First(&stockBatch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundEntity(stockbatchmodel.EntityName, err)
		}
		return nil, err
	}
	return &stockBatch, nil
}

// GetStockBatchesByIngredientIDs fetches all stock batches for given ingredientIDs ordered by expiration_date ASC.
func (s *mysqlStockBatch) GetStockBatchesByIngredientIDs(ctx context.Context,
	ingredientIDs []uint64) ([]stockbatchmodel.StockBatch, error) {
	if len(ingredientIDs) == 0 {
		return nil, fmt.Errorf("no ingredient IDs provided")
	}

	var stockBatches []stockbatchmodel.StockBatch
	err := s.db.WithContext(ctx).
		Where("ingredient_id IN ?", ingredientIDs).
		Order("ingredient_id ASC, expiration_date ASC").
		Find(&stockBatches).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock batches: %w", err)
	}

	return stockBatches, nil
}
