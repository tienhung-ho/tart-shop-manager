package rolestorage

import (
	"context"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

func (s *mysqlRole) DeleteRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db

	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	var role rolemodel.Role
	if err := db.WithContext(ctx).Where(cond).Delete(&role).Error; err != nil {
		return err
	}
	return nil
}
