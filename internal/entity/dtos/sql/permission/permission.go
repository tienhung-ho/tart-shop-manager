package permissionmodel

import "tart-shop-manager/internal/common"

const (
	EntityName = "permission"
)

type Role struct {
	Id          int           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string        `json:"name" gorm:"column:name;not null;unique"`
	Description string        `json:"description" gorm:"column:description;type:text"`
	Status      common.Status `json:"status" gorm:"column:status;type:enum('Active','Inactive','Pending')"`
}

type Permission struct {
	PermissionID uint   `gorm:"column:permission_id;primaryKey;autoIncrement" json:"permission_id"`
	Name         string `gorm:"column:name;size:255;not null;unique" json:"name"`
	Description  string `gorm:"column:description" json:"description"`
	Roles        []Role `gorm:"many2many:role_permission;"`
	*common.CommonFields
}

func (Permission) TableName() string {
	return "Permission"
}
