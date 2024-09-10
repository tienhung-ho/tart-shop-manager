package productstorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (s *mysqlProduct) DeleteProduct(ctx context.Context, cond map[string]interface{}, morekyes ...string) error {
	db := s.db

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	var product productmodel.Product
	if err := db.WithContext(ctx).Where(cond).Delete(&product).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
