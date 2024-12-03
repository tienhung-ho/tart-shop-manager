package productstorage

import "gorm.io/gorm"

var AllowedSortFields = map[string]bool{
	"name":       true,
	"created_at": true,
	"updated_at": true,
	"product_id": true,
}

type mysqlProduct struct {
	db *gorm.DB
}

func NewMySQLProduct(db *gorm.DB) *mysqlProduct {
	return &mysqlProduct{db}
}
