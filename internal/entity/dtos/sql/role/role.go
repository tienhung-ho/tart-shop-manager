package rolemodel

import (
	"tart-shop-manager/internal/common"
	rolecachemodel "tart-shop-manager/internal/entity/dtos/redis/role"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

const (
	EntityName = "role"
)

type Role struct {
	RoleID      uint                         `gorm:"column:role_id;primaryKey;autoIncrement" json:"role_id"`
	Name        string                       `gorm:"column:name;size:255;not null;unique" json:"name"`
	Description string                       `gorm:"column:description" json:"description"`
	Permissions []permissionmodel.Permission `gorm:"many2many:role_permission;"`
	*common.CommonFields
}

func (Role) TableName() string {
	return "Role"
}

func (r Role) ToCreateRoleCache() *rolecachemodel.CreateRole {
	return &rolecachemodel.CreateRole{
		RoleID:      r.RoleID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		CommonFields: &common.CommonFields{
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Status:    r.Status,
			CreatedBy: r.CreatedBy,
			UpdatedBy: r.UpdatedBy,
		},
	}
}