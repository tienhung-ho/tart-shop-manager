package imagestorage

import (
	"context"
	"log"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

func (s *mysqlImage) UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error {

	db := s.db.Begin()

	log.Print(cond, data.ProductID)

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Where(cond).Updates(data).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
