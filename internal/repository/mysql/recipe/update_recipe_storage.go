package recipestorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlRecipe) UpdateRecipe(ctx context.Context, cond map[string]interface{},
	data *recipemodel.UpdateRecipe, morekeys ...string) (*recipemodel.Recipe, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&recipemodel.UpdateRecipe{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, recipemodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(recipemodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return nil, err
	}

	var record recipemodel.Recipe

	if err := db.WithContext(ctx).Model(data).Select(recipemodel.SelectFields).
		Where(cond).Preload("Product").First(&record).Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
