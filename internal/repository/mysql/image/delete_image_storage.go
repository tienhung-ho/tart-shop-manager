package imagestorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

func (s *mysqlImage) DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db.Begin()

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	var image imagemodel.Image
	if err := db.WithContext(ctx).Where(cond).Delete(&image).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
