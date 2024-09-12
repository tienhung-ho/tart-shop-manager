package categorybusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	responseutil "tart-shop-manager/internal/util/response"
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

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, categorymodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(categorymodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	return recordID, nil
}
