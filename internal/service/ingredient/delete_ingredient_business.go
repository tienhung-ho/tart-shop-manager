package ingredientbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteIngredientStorage interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
	DeleteIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteIngredientCache interface {
	DeleteIngredient(ctx context.Context, morekeys ...string) error
}

type deleteIngredientBusiness struct {
	store DeleteIngredientStorage
	cache DeleteIngredientCache
}

func NewDeleteIngredientBiz(store DeleteIngredientStorage, cache DeleteIngredientCache) *deleteIngredientBusiness {
	return &deleteIngredientBusiness{
		store: store,
		cache: cache,
	}
}

func (biz *deleteIngredientBusiness) DeleteIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetIngredient(ctx, cond, morekeys...)

	if err != nil {
		return common.ErrNotFoundEntity(ingredientmodel.EntityName, err)
	}

	if record == nil {
		return common.ErrRecordExist(ingredientmodel.EntityName, nil)
	}

	if err := biz.store.DeleteIngredient(ctx, cond, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(ingredientmodel.EntityName, err)
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

	if err := biz.cache.DeleteIngredient(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(ingredientmodel.EntityName, err)
	}

	return nil
}
