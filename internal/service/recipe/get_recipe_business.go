package recipebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetRecipeStorage interface {
	GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
}

type GetRecipeCache interface {
	GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
	SaveRecipe(ctx context.Context, data interface{}, morekeys ...string) error
}

type getRecipeBusiness struct {
	store GetRecipeStorage
	cache GetRecipeCache
}

func NewGetRecipeBiz(store GetRecipeStorage, cache GetRecipeCache) *getRecipeBusiness {
	return &getRecipeBusiness{store: store, cache: cache}
}

func (biz *getRecipeBusiness) GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error) {

	record, err := biz.cache.GetRecipe(ctx, cond, morekeys...)
	if err != nil {
		return nil, common.ErrNotFoundEntity(recipemodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetRecipe(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(recipemodel.EntityName, err)
	}

	if record != nil {
		var pagging paggingcommon.Paging
		pagging.Process()

		var createRecipe = record.ToCreateRecipe()

		// Generate cache key
		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: recipemodel.EntityName,
			Cond:       cond,
			Paging:     pagging,
			Filter:     commonfilter.Filter{},
			MoreKeys:   morekeys,
		})
		if err != nil {
			return nil, common.ErrCannotGenerateKey(recipemodel.EntityName, err)
		}

		if err := biz.cache.SaveRecipe(ctx, createRecipe, key); err != nil {
			return nil, common.ErrCannotCreateEntity(recipemodel.EntityName, err)
		}
	}

	return record, nil
}
