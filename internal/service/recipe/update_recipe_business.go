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

type UpdateRecipeStorage interface {
	GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
	UpdateRecipe(ctx context.Context, cond map[string]interface{}, data *recipemodel.UpdateRecipe, morekeys ...string) (*recipemodel.Recipe, error)
}

type UpdateRecipeCache interface {
	DeleteRecipe(ctx context.Context, morekeys ...string) error
}

type updateRecipeBusiness struct {
	store UpdateRecipeStorage
	cache UpdateRecipeCache
}

func NewUpdateRecipeBiz(store UpdateRecipeStorage, cache UpdateRecipeCache) *updateRecipeBusiness {
	return &updateRecipeBusiness{store, cache}
}

func (biz *updateRecipeBusiness) UpdateRecipe(ctx context.Context,
	cond map[string]interface{}, data *recipemodel.UpdateRecipe, morekeys ...string) (*recipemodel.Recipe, error) {

	record, err := biz.store.GetRecipe(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrCannotGetEntity(recipemodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(recipemodel.EntityName, err)
	}

	updateRecord, err := biz.store.UpdateRecipe(ctx, map[string]interface{}{"recipe_id": record.RecipeID}, data, morekeys...)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(recipemodel.EntityName, err)
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
		return nil, common.ErrCannotGenerateKey(recipemodel.EntityName, err)
	}

	if err := biz.cache.DeleteRecipe(ctx, key); err != nil {
		return nil, common.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}
	return updateRecord, nil
}
