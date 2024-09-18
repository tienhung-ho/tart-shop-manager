package stockbatchstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (s *mysqlStockBatch) CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		db.Rollback()
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.StockBatchID, nil
}
