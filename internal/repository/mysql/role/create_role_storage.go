package rolestorage

import (
	"context"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

func (s *mysqlRole) CreateRole(ctx context.Context, data *rolemodel.CreateRole, morekeys ...string) (uint, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	// Link role with permissions if any permissions are provided
	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	// Link role with permissions if any permissions are provided
	if len(data.Permissions) > 0 {
		if err := db.Model(data).Association("Permissions").Replace(data.Permissions); err != nil {
			db.Rollback()
			return 0, common.ErrDB(err)
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.RoleID, nil
}
