package stockbatchbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

type CreateStockBatchStorage interface {
	CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint, error)
}

type createStockBatchBusiness struct {
	store CreateStockBatchStorage
}

func NewCreateStockBatchBiz(store CreateStockBatchStorage) *createStockBatchBusiness {
	return &createStockBatchBusiness{store: store}
}

func (biz *createStockBatchBusiness) CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint, error) {

	recordID, err := biz.store.CreateStockBatch(ctx, data, morekeys...)

	if err != nil {
		return 0, common.ErrCannotUpdateEntity(stockbatchmodel.EntityName, err)
	}

	return recordID, nil
}
