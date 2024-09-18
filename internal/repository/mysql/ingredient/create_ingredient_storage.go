package ingredientstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) CreateIngredient(ctx context.Context, ingredient *ingredientmodel.CreateIngredient, morekeys ...string) (uint, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(&ingredient).Error; err != nil {
		db.Rollback()
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return ingredient.IngredientID, nil
}
