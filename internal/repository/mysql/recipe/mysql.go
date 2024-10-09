package recipestorage

import "gorm.io/gorm"

var (
	SelectFields = []string{
		"recipe_id",
		"size",
		"description",
		"cost",
		"status",
	}
)

type mysqlRecipe struct {
	db *gorm.DB
}

func NewMySQLRecipe(db *gorm.DB) *mysqlRecipe {
	return &mysqlRecipe{db}
}
