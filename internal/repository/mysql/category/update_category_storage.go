package categorystorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlCategory) UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	// Sử dụng con trỏ đến struct thực tế
	if err := db.WithContext(ctx).Model(&categorymodel.UpdateCategory{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return nil, err
	}

	var updatedCategory categorymodel.Category
	if err := db.WithContext(ctx).Select(categorymodel.SelectFields).Where(cond).Preload("Images").First(&updatedCategory).Error; err != nil {
		db.Rollback()
		return nil, common.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &updatedCategory, nil
}
