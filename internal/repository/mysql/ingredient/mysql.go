package ingredientstorage

import "gorm.io/gorm"

type mysqlIngredient struct {
	db *gorm.DB
}

func NewMySQLIngredient(db *gorm.DB) *mysqlIngredient {
	return &mysqlIngredient{db}
}
