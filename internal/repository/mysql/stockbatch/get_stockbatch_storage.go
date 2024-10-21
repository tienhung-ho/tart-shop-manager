package stockbatchstorage

import (
	"context"
	"errors"
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
