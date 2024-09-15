package categorybusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
	responseutil "tart-shop-manager/internal/util/response"
)

type UpdateCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error)
}

type UpdateCategoryCache interface {
	DeleteCategory(ctx context.Context, morekeys ...string) error
}

type updateCategoryBusiness struct {
	store UpdateCategoryStorage
	cache UpdateCategoryCache
}

func NewUpdateCategoryBiz(store UpdateCategoryStorage, cache UpdateCategoryCache) *updateCategoryBusiness {
	return &updateCategoryBusiness{store, cache}
}

func (biz *updateCategoryBusiness) UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error) {

	record, err := biz.store.GetCategory(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateCategory(ctx, map[string]interface{}{"category_id": record.CategoryID}, data)
	if err != nil {
		// Check for MySQL duplicate entry error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, categorymodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(categorymodel.EntityName, fieldName, err)
		}

		return nil, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	key := cacheutil.GenerateKey(accountmodel.EntityName, cond, pagging, commonfilter.Filter{})

	if err := biz.cache.DeleteCategory(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	return updatedRecord, nil
}
