package recipebusiness

import (
	"context"
	"fmt"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
	databaseutil "tart-shop-manager/internal/util/database"
)

type CreateRecipeStorage interface {
	CreateRecipe(ctx context.Context, data *recipemodel.CreateRecipe) (uint64, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type CreateRecipeCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type createRecipeBusiness struct {
	store                 CreateRecipeStorage
	cache                 CreateRecipeCache
	ingredientStore       ListItemIngredientStorage
	recipeIngredientStore RecipeIngredientStorage
}

func NewCreateRecipeBusiness(store CreateRecipeStorage, cache CreateRecipeCache, ingredientStore ListItemIngredientStorage, recipeIngredientStore RecipeIngredientStorage) *createRecipeBusiness {
	return &createRecipeBusiness{store, cache, ingredientStore, recipeIngredientStore}
}

func (biz *createRecipeBusiness) CreateRecipe(ctx context.Context, data *recipemodel.CreateRecipe) (uint64, error) {

	// Lấy danh sách ingredient_id từ request
	ingredientIDs := make([]uint64, len(data.Ingredients))
	for i, ing := range data.Ingredients {
		ingredientIDs[i] = ing.IngredientID
	}

	// Sử dụng hàm ListItem để lấy danh sách nguyên liệu tồn tại
	cond := map[string]interface{}{}
	paging := &paggingcommon.Paging{
		Page:  1,
		Limit: len(ingredientIDs),
	}
	filter := &commonfilter.Filter{
		IDs: ingredientIDs,
	}

	existingIngredients, err := biz.ingredientStore.ListItem(ctx, cond, paging, filter)
	if err != nil {
		return 0, common.ErrCannotGetEntity("ingredients", err)
	}

	// So sánh danh sách ingredientIDs và existingIngredients để tìm ra các ingredient_id không tồn tại
	existingIngredientIDs := make([]uint64, len(existingIngredients))
	for i, ing := range existingIngredients {
		existingIngredientIDs[i] = uint64(ing.IngredientID)
	}

	missingIngredientIDs := databaseutil.Difference(ingredientIDs, existingIngredientIDs)
	if len(missingIngredientIDs) > 0 {
		return 0, common.ErrInvalidRequest(fmt.Errorf("ingredients not found: %v", missingIngredientIDs))
	}
	var recordID uint64
	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {
		recordID, err = biz.store.CreateRecipe(ctx, data)

		if err != nil {
			return common.ErrCannotCreateEntity(recipemodel.EntityName, err)
		}

		var createRecipeIngredient []recipeingredientmodel.RecipeIngredientCreate
		for _, ing := range data.Ingredients {
			createRecipeIngredient = append(createRecipeIngredient, recipeingredientmodel.RecipeIngredientCreate{
				RecipeID:     data.RecipeID,
				IngredientID: ing.IngredientID,
				Quantity:     ing.Quantity,
				Unit:         ing.Unit,
			})
		}
		if len(createRecipeIngredient) > 0 {
			err = biz.recipeIngredientStore.CreateRecipeIngredients(txCtx, createRecipeIngredient)
			if err != nil {
				return common.ErrCannotCreateEntity("recipe ingredients", err)
			}
		}

		if err := biz.cache.DeleteListCache(txCtx, recipemodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(recipemodel.EntityName, err)
		}

		if err := biz.cache.DeleteListCache(ctx, productmodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(recipemodel.EntityName, err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return recordID, nil
}
