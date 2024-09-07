package rolestorage

import (
	"context"
	commonrecover "tart-shop-manager/internal/common/recover"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

func (s *mysqlRole) GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error) {

	db := s.db

	defer commonrecover.RecoverTransaction(db)

	var role rolemodel.Role

	if err := db.WithContext(ctx).Where(cond).Preload("Permissions").First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
