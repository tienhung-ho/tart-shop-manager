package rolecachemodel

import (
	"tart-shop-manager/internal/common"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

type CreateRole struct {
	RoleID      uint                         `json:"role_id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Permissions []permissionmodel.Permission `gorm:"many2many:role_permission;"`
	*common.CommonFields
}
