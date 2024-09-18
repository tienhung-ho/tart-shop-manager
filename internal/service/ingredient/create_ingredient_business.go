package ingredientbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	responseutil "tart-shop-manager/internal/util/response"
)

type CreateIngredientStorage interface {
	GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error)
	CreateIngredient(ctx context.Context, ingredient *ingredientmodel.CreateIngredient, morekeys ...string) (uint, error)
}

type createIngredientBusiness struct {
	store CreateIngredientStorage
}

func NewCreateIngredientBiz(store CreateIngredientStorage) *createIngredientBusiness {
	return &createIngredientBusiness{store: store}
}

func (biz *createIngredientBusiness) CreateIngredient(ctx context.Context, ingredient *ingredientmodel.CreateIngredient, morekeys ...string) (uint, error) {

	// Kiểm tra xem ingredient đã tồn tại hay chưa
	record, err := biz.store.GetIngredient(ctx, map[string]interface{}{"name": ingredient.Name})

	if err != nil {
		return 0, common.ErrCannotGetEntity(ingredientmodel.EntityName, err)
	}

	if record != nil {
		return 0, common.ErrRecordExist(ingredientmodel.EntityName, nil)
	}

	recordId, err := biz.store.CreateIngredient(ctx, ingredient, morekeys...)
	if err != nil {
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(ingredientmodel.EntityName, err)
	}

	return recordId, nil
}
