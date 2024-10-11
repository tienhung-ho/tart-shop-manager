package recipebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteRecipeStorage interface {
	GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
	DeleteRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteRecipeCache interface {
	DeleteRecipe(ctx context.Context, morekeys ...string) error
}

type deleteRecipeBusiness struct {
	store DeleteRecipeStorage
	cache DeleteRecipeCache
}

func NewDeleteRecipeBiz(store DeleteRecipeStorage, cache DeleteRecipeCache) *deleteRecipeBusiness {
	return &deleteRecipeBusiness{store, cache}
}

func (biz *deleteRecipeBusiness) DeleteRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetRecipe(ctx, cond, morekeys...)

	if err != nil {
		return common.ErrCannotGetEntity(recipemodel.EntityName, err)
	}

	if record == nil {
		return common.ErrNotFoundEntity(recipemodel.EntityName, err)
	}

	if err := biz.store.DeleteRecipe(ctx, map[string]interface{}{"recipe_id": record.RecipeID}, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(recipemodel.EntityName, err)
	}

	// 6. Xóa cache sản phẩm
	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: recipemodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return common.ErrCannotGenerateKey(recipemodel.EntityName, err)
	}

	if err := biz.cache.DeleteRecipe(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	return nil
}
