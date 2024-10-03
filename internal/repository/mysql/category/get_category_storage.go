package categorystorage

import (
	"context"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

func (s *mysqlCategory) GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {

	db := s.db

	var record categorymodel.Category
	if err := db.WithContext(ctx).Select(categorymodel.SelectFields).
		Where(cond).Preload("Images").First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
