package categorybusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteCategoryCache interface {
	DeleteCategory(ctx context.Context, morekeys ...string) error
}

type deleteCategoryBusiness struct {
	store DeleteCategoryStorage
	cache DeleteCategoryCache
}

func NewDeleteCategoryBiz(store DeleteCategoryStorage, cache DeleteCategoryCache) *deleteCategoryBusiness {
	return &deleteCategoryBusiness{store, cache}
}

func (biz *deleteCategoryBusiness) DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetCategory(ctx, cond)

	if err != nil {
		return common.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	if err := biz.store.DeleteCategory(ctx, map[string]interface{}{"category_id": record.CategoryID}, morekeys...); err != nil {
		return common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: categorymodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return common.ErrCannotGenerateKey(categorymodel.EntityName, err)
	}

	if err := biz.cache.DeleteCategory(ctx, key); err != nil {
		return common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	return nil
}
