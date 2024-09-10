package productbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	responseutil "tart-shop-manager/internal/util/response"
)

type CreateProductStorage interface {
	GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
	CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error)
}

type createProductBusiness struct {
	store CreateProductStorage
}

func NewCreateProductBusiness(store CreateProductStorage) *createProductBusiness {
	return &createProductBusiness{store}
}

func (biz *createProductBusiness) CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error) {

	recordId, err := biz.store.CreateProduct(ctx, data)

	if err != nil {
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	return recordId, nil
}
