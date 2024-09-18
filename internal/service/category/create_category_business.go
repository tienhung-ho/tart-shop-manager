package categorybusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

type CreateCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	CreateCategory(ctx context.Context, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error)
}

type createCategoryBusiness struct {
	store CreateCategoryStorage
}

func NewCreateCategoryBusiness(store CreateCategoryStorage) *createCategoryBusiness {
	return &createCategoryBusiness{store}
}

func (biz *createCategoryBusiness) CreateCategory(ctx context.Context, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error) {

	recordID, err := biz.store.CreateCategory(ctx, data)

	if err != nil {
		// Check for MySQL duplicate entry error
		return 0, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	return recordID, nil
}
