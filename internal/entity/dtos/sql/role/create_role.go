package rolemodel

import (
	"tart-shop-manager/internal/common"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

type CreateRole struct {
	RoleID      uint                         `gorm:"column:role_id;primaryKey;autoIncrement" json:"-"`
	Name        string                       `gorm:"column:name;size:255;not null;unique" json:"name"`
	Description string                       `gorm:"column:description" json:"description"`
	Permissions []permissionmodel.Permission `gorm:"many2many:role_permissions;foreignKey:RoleID;joinForeignKey:RoleID;References:PermissionID;joinReferences:PermissionID"`
	common.CommonFields
}

func (CreateRole) TableName() string {
	return Role{}.TableName()
}
