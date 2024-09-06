package permissionmodel

import "tart-shop-manager/internal/common"

const (
	EntityName = "permission"
)

type Permission struct {
	PermissionID uint   `gorm:"column:permission_id;primaryKey;autoIncrement" json:"permission_id"`
	Name         string `gorm:"column:name;size:255;not null;unique" json:"name"`
	Object       string `gorm:"column:object;size:255;not null" json:"object"` // Object được truy cập
	Action       string `gorm:"column:action;size:50;not null" json:"action"`  // Hành động: POST, CREATE, DELETE...
	Description  string `gorm:"column:description" json:"description"`
	*common.CommonFields
}

func (Permission) TableName() string {
	return "Permission"
}
