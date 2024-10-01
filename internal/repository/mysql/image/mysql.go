package imagestorage

import "gorm.io/gorm"

type mysqlImage struct {
	db *gorm.DB
}

func NewMySQLImage(db *gorm.DB) *mysqlImage {
	return &mysqlImage{db}
}
