package recipestorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlRecipe) CreateRecipe(ctx context.Context, data *recipemodel.CreateRecipe) (uint64, error) {
	db := s.getDB(ctx)

	if err := db.Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := responseutil.ExtractFieldFromError(err, recipemodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(recipemodel.EntityName, fieldName, err)
		}
		return 0, common.ErrDB(err)
	}

	return data.RecipeID, nil
}
