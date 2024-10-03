package imagestorage

import (
	"context"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

func (s *mysqlImage) ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error) {

	db := s.db

	var records []imagemodel.Image
	if err := db.WithContext(ctx).Where(cond).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}
