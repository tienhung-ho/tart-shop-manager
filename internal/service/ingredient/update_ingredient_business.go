package ingredientbusiness

import (
	"context"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateIngredientStorage interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
	UpdateIngredient(ctx context.Context, cond map[string]interface{},
		ingredient *ingredientmodel.UpdateIngredient, morekeys ...string) (*ingredientmodel.Ingredient, error)
}

type UpdateIngredientCache interface {
	DeleteIngredient(ctx context.Context, morekeys ...string) error
}

type updateIngredientBusiness struct {
	store UpdateIngredientStorage
	cache UpdateIngredientCache
}

func NewUpdateIngredientBiz(store UpdateIngredientStorage, cache UpdateIngredientCache) *updateIngredientBusiness {
	return &updateIngredientBusiness{store, cache}
}

func (biz *updateIngredientBusiness) UpdateIngredient(ctx context.Context, cond map[string]interface{},
	ingredient *ingredientmodel.UpdateIngredient, morekeys ...string) (*ingredientmodel.Ingredient, error) {

	record, err := biz.store.GetIngredient(ctx, cond)

	if err != nil {
		return nil, common.ErrNotFoundEntity(ingredientmodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrRecordExist(ingredientmodel.EntityName, nil)
	}

	updatedRecord, err := biz.store.UpdateIngredient(ctx, cond, ingredient, morekeys...)

	if err != nil {
		log.Print(err)
		return nil, common.ErrCannotUpdateEntity(ingredientmodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: ingredientmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(ingredientmodel.EntityName, err)
	}

	if err := biz.cache.DeleteIngredient(ctx, key); err != nil {
		log.Print(err)
		return nil, common.ErrCannotDeleteEntity(ingredientmodel.EntityName, err)
	}

	return updatedRecord, nil
}
