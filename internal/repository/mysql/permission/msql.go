package permissionstorage

import "gorm.io/gorm"

type mysqlPermission struct {
	db *gorm.DB
}

func NewMySQLPermission(db *gorm.DB) *mysqlPermission {
	return &mysqlPermission{db}
}
