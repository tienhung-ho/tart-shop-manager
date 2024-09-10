package productstorage

import "gorm.io/gorm"

type mysqlProduct struct {
	db *gorm.DB
}

func NewMySQLProduct(db *gorm.DB) *mysqlProduct {
	return &mysqlProduct{db}
}
