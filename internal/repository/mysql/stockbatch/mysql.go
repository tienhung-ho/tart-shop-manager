package rolestorage

import "gorm.io/gorm"

type mysqlRole struct {
	db *gorm.DB
}

func NewMySQLRole(db *gorm.DB) *mysqlRole {
	return &mysqlRole{db}
}
