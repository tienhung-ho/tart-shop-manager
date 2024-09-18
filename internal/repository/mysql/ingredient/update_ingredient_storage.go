package ingredientstorage

import (
	"context"
	"log"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) UpdateIngredient(ctx context.Context, cond map[string]interface{},
	data *ingredientmodel.UpdateIngredient, morekeys ...string) (*ingredientmodel.Ingredient, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&ingredientmodel.UpdateIngredient{}).Where(cond).
		//Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(&data).Error; err != nil {

		//var mysqlErr *mysql.MySQLError
		//if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		//
		//	fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
		//	return nil, common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		//}
		log.Print(err)
		db.Rollback()
		return nil, err

	}

	var record ingredientmodel.Ingredient

	if err := db.WithContext(ctx).Model(data).Where(cond).First(&record).Error; err != nil {
		log.Print(err)
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		log.Print(err)
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
