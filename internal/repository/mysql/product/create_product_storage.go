package productstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (s *mysqlProduct) CreateProduct(ctx context.Context, data *productmodel.CreateProduct, morekeys ...string) (uint64, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		return 0, err
	}

	if err := db.Commit().Error; err != nil {
		return 0, common.ErrDB(err)
	}

	return data.ProductID, nil
}
