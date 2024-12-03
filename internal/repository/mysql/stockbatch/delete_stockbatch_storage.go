package stockbatchstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (r *mysqlStockBatch) DeleteStockBatches(ctx context.Context, stockBatchIDs []uint64) error {
	if len(stockBatchIDs) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).
		Where("stockbatch_id IN ?", stockBatchIDs).
		Delete(&stockbatchmodel.StockBatch{}).Error; err != nil {
		return common.ErrCannotDeleteEntity("StockBatch", err)
	}
	return nil
}

func (r *mysqlStockBatch) DeleteStockBatch(ctx context.Context, cond map[string]interface{}) error {

	if err := r.db.WithContext(ctx).
		Where(cond).
		Delete(&stockbatchmodel.StockBatch{}).Error; err != nil {
		return common.ErrCannotDeleteEntity(stockbatchmodel.EntityName, err)
	}

	return nil
}
