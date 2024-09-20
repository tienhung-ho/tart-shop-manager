package categorybusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
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
		return nil, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: categorymodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(categorymodel.EntityName, err)
	}

	if err := biz.cache.DeleteCategory(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	return updatedRecord, nil
}
