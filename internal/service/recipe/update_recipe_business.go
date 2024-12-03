package recipebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
	cacheutil "tart-shop-manager/internal/util/cache"
	databaseutil "tart-shop-manager/internal/util/database"
)

type UpdateRecipeStorage interface {
	GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error)
	UpdateRecipe(ctx context.Context, cond map[string]interface{}, data *recipemodel.UpdateRecipe, morekeys ...string) (*recipemodel.Recipe, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UpdateRecipeCache interface {
	DeleteRecipe(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type updateRecipeBusiness struct {
	store                 UpdateRecipeStorage
	cache                 UpdateRecipeCache
	recipeIngredientStore RecipeIngredientStorage
}

func NewUpdateRecipeBiz(store UpdateRecipeStorage, cache UpdateRecipeCache, recipeIngredientStore RecipeIngredientStorage) *updateRecipeBusiness {
	return &updateRecipeBusiness{store,
		cache,
		recipeIngredientStore}
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

	var updatedRecord *recipemodel.Recipe

	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {
		// Cập nhật Recipe trước
		updatedRecord, err = biz.store.UpdateRecipe(txCtx, map[string]interface{}{"recipe_id": record.RecipeID}, data, morekeys...)
		if err != nil {
			return common.ErrCannotUpdateEntity(recipemodel.EntityName, err)
		}

		// Kiểm tra nếu RecipeIngredients không rỗng thì mới thực hiện cập nhật
		if data.RecipeIngredients != nil && len(data.RecipeIngredients) > 0 {
			existingIngredients := record.RecipeIngredients

			// Lấy danh sách ingredient_id hiện có
			existingIngredientIDs := make([]uint64, len(existingIngredients))
			for i, ing := range existingIngredients {
				existingIngredientIDs[i] = ing.IngredientID
			}

			// Từ request
			newIngredients := data.RecipeIngredients // []RecipeIngredientCreate

			newIngredientIDs := make([]uint64, len(newIngredients))
			newIngredientsMap := make(map[uint64]recipeingredientmodel.RecipeIngredientCreate)
			for i, ing := range newIngredients {
				newIngredientIDs[i] = ing.IngredientID
				newIngredientsMap[ing.IngredientID] = recipeingredientmodel.RecipeIngredientCreate{
					RecipeID:     record.RecipeID,
					Unit:         ing.Unit,
					IngredientID: ing.IngredientID,
					Quantity:     ing.Quantity,
				}
			}

			// Ingredients cần thêm mới
			ingredientIDsToAdd := databaseutil.Difference(newIngredientIDs, existingIngredientIDs)

			// Ingredients cần xóa
			ingredientIDsToRemove := databaseutil.Difference(existingIngredientIDs, newIngredientIDs)

			// Ingredients cần cập nhật
			var ingredientsToUpdate []recipeingredientmodel.RecipeIngredient
			for _, existingIngredient := range existingIngredients {
				if newIng, ok := newIngredientsMap[existingIngredient.IngredientID]; ok {
					if existingIngredient.Quantity != newIng.Quantity || existingIngredient.Unit != newIng.Unit {
						ingredientsToUpdate = append(ingredientsToUpdate, recipeingredientmodel.RecipeIngredient{
							RecipeID:     existingIngredient.RecipeID,
							IngredientID: existingIngredient.IngredientID,
							Quantity:     newIng.Quantity,
							Unit:         newIng.Unit,
						})
					}
				}
			}

			// Thực hiện thêm mới ingredients
			var ingredientsToAddModels []recipeingredientmodel.RecipeIngredientCreate
			for _, ingID := range ingredientIDsToAdd {
				newIng := newIngredientsMap[ingID]
				ingredientsToAddModels = append(ingredientsToAddModels, recipeingredientmodel.RecipeIngredientCreate{
					RecipeID:     record.RecipeID,
					IngredientID: ingID,
					Quantity:     newIng.Quantity,
					Unit:         newIng.Unit,
				})
			}

			if len(ingredientsToAddModels) > 0 {
				err = biz.recipeIngredientStore.CreateRecipeIngredients(txCtx, ingredientsToAddModels)
				if err != nil {
					return common.ErrCannotCreateEntity("recipe ingredients", err)
				}
			}

			// Thực hiện xóa ingredients
			if len(ingredientIDsToRemove) > 0 {
				err := biz.recipeIngredientStore.RemoveRecipeIngredients(txCtx, record.RecipeID, ingredientIDsToRemove)
				if err != nil {
					return common.ErrCannotDeleteEntity("recipe ingredients", err)
				}
			}

			// Thực hiện cập nhật ingredients
			if len(ingredientsToUpdate) > 0 {
				for _, ing := range ingredientsToUpdate {
					updateIng := ing.ToUpdateRecipeIngredient()
					err := biz.recipeIngredientStore.UpdateRecipeIngredient(txCtx, updateIng)
					if err != nil {
						return common.ErrCannotUpdateEntity("recipe ingredient", err)
					}
				}
			}
		}

		return nil
	})

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
	if err := biz.cache.DeleteListCache(ctx, recipemodel.EntityName); err != nil {
		return nil, common.ErrCannotDeleteEntity(recipemodel.EntityName, err)
	}

	if err := biz.cache.DeleteListCache(ctx, productmodel.EntityName); err != nil {
		return nil, common.ErrCannotDeleteEntity(recipemodel.EntityName, err)
	}

	return updatedRecord, nil
}
