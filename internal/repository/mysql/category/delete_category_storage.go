package categorystorage

import (
	"context"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

func (s *mysqlCategory) DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db

	var record categorymodel.Category
	if err := db.WithContext(ctx).Where(cond).Delete(&record).Error; err != nil {
		return err
	}

	return nil
}
