package recipestorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

func (s *mysqlRecipe) DeleteRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db.Begin()

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	var record recipemodel.Recipe
	if err := db.WithContext(ctx).Where(cond).Delete(&record).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
