package stockbatchstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (s *mysqlStockBatch) GetStockBatch(ctx context.Context, stockBatchID uint64) (*stockbatchmodel.StockBatch, error) {
	var stockBatch stockbatchmodel.StockBatch
	err := s.db.WithContext(ctx).
		Where("stockbatch_id = ?", stockBatchID).
		First(&stockBatch).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundEntity(stockbatchmodel.EntityName, err)
		}
		return nil, err
	}
	return &stockBatch, nil
}
