package ingredientbusiness

import (
	"context"
	"errors"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

type CreateIngredientStorage interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
	CreateIngredient(ctx context.Context, ingredient *ingredientmodel.CreateIngredient, morekeys ...string) (uint64, error)
}

type CreateIngredientCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type createIngredientBusiness struct {
	store CreateIngredientStorage
	cache CreateIngredientCache
}

func NewCreateIngredientBiz(store CreateIngredientStorage, cache CreateIngredientCache) *createIngredientBusiness {
	return &createIngredientBusiness{store: store, cache: cache}
}

func (biz *createIngredientBusiness) CreateIngredient(ctx context.Context, ingredient *ingredientmodel.CreateIngredient, morekeys ...string) (uint64, error) {

	// Kiểm tra xem ingredient đã tồn tại hay chưa
	record, err := biz.store.GetIngredient(ctx, map[string]interface{}{"name": ingredient.Name})

	if err != nil && !errors.Is(err, common.RecordNotFound) {
		return 0, common.ErrCannotGetEntity(ingredientmodel.EntityName, err)
	}

	if record != nil {
		return 0, common.ErrRecordExist(ingredientmodel.EntityName, nil)
	}

	recordId, err := biz.store.CreateIngredient(ctx, ingredient, morekeys...)
	if err != nil {
		return 0, common.ErrCannotCreateEntity(ingredientmodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, ingredientmodel.EntityName); err != nil {
		return 0, common.ErrCannotDeleteEntity(ingredientmodel.EntityName, err)
	}

	return recordId, nil
}
