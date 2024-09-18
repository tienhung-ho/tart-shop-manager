package stockbatchbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	responseutil "tart-shop-manager/internal/util/response"
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
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	return recordID, nil
}
