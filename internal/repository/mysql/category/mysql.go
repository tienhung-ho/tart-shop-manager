package categorystorage

import "gorm.io/gorm"

type mysqlCategory struct {
	db *gorm.DB
}

func NewMySQLCategory(db *gorm.DB) *mysqlCategory {
	return &mysqlCategory{db}
}
