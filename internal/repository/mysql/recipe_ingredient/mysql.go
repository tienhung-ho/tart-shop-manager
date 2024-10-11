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
	AllowedSortFields = map[string]bool{
		"recipe_id":  true,
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"cost":       true,
	}
)

type mysqlRecipe struct {
	db *gorm.DB
}

func NewMySQLRecipe(db *gorm.DB) *mysqlRecipe {
	return &mysqlRecipe{db}
}
