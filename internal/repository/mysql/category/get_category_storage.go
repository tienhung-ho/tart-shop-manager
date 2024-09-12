package categorystorage

import (
	"context"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

func (s *mysqlCategory) GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {

	db := s.db

	var record categorymodel.Category
	if err := db.WithContext(ctx).Where(cond).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
