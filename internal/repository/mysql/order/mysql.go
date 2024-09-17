package orderstorage

import "gorm.io/gorm"

type mysqlOrder struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlOrder {
	return &mysqlOrder{db}
}
