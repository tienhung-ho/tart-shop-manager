package ingredientbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetIngredientStorage interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
}

type GetIngredientCache interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
	SaveIngredient(ctx context.Context, data interface{}, morekeys ...string) error
}

type getIngredientBusiness struct {
	store GetIngredientStorage
	cache GetIngredientCache
}

func NewGetIngredientBiz(store GetIngredientStorage, cache GetIngredientCache) *getIngredientBusiness {
	return &getIngredientBusiness{store: store, cache: cache}
}

func (biz *getIngredientBusiness) GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error) {

	record, err := biz.cache.GetIngredient(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(ingredientmodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetIngredient(ctx, cond, morekeys...)
	if err != nil {
		return nil, common.ErrNotFoundEntity(ingredientmodel.EntityName, err)
	}

	if record != nil {
		var pagging paggingcommon.Paging
		pagging.Process()

		key := cacheutil.GenerateKey(ingredientmodel.EntityName, cond, pagging, commonfilter.Filter{})

		cacheRecord := record.ToCreateIngredientCache()

		if err := biz.cache.SaveIngredient(ctx, cacheRecord, key); err != nil {
			return nil, common.ErrCannotCreateEntity(ingredientmodel.EntityName, err)
		}
	}

	return record, nil
}
