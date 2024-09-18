package productbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
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
		return 0, common.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	return recordId, nil
}
