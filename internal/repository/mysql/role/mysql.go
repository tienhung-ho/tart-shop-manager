package rolestorage

import "gorm.io/gorm"

var (
	AllowedSortFields = map[string]bool{
		"role_id":    true,
		"created_at": true,
		"updated_at": true,
	}
)

type mysqlRole struct {
	db *gorm.DB
}

func NewMySQLRole(db *gorm.DB) *mysqlRole {
	return &mysqlRole{db}
}
