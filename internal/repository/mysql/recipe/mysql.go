package recipestorage

import "gorm.io/gorm"

type mysqlRecipe struct {
	db *gorm.DB
}

func NewMySQLRecipe(db *gorm.DB) *mysqlRecipe {
	return &mysqlRecipe{db}
}
