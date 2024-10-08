package recipebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

type CreateRecipeStorage interface {
	CreateRecipe(ctx context.Context, data *recipemodel.CreateRecipe) (uint64, error)
}

type createRecipeBusiness struct {
	store CreateRecipeStorage
}

func NewCreateRecipeBusiness(store CreateRecipeStorage) *createRecipeBusiness {
	return &createRecipeBusiness{store}
}

func (biz *createRecipeBusiness) CreateRecipe(ctx context.Context, data *recipemodel.CreateRecipe) (uint64, error) {

	recordId, err := biz.store.CreateRecipe(ctx, data)

	if err != nil {
		return 0, common.ErrCannotCreateEntity(recipemodel.EntityName, err)
	}

	return recordId, nil
}
