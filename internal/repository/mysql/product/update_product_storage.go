package productstorage

import (
	"context"
	"gorm.io/gorm/clause"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (s *mysqlProduct) UpdateProduct(ctx context.Context, cond map[string]interface{}, data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&productmodel.UpdateProduct{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(data).Error; err != nil {

		db.Rollback()
		return nil, err
	}

	var record productmodel.Product

	if err := db.WithContext(ctx).Model(data).Where(cond).Preload("Category").Preload("Recipes").First(&record).Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
