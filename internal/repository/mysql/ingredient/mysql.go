package ingredientstorage

import "gorm.io/gorm"

var (
	AllowedSortFields = map[string]bool{
		"name":          true,
		"created_at":    true,
		"updated_at":    true,
		"ingredient_id": true,
	}
)

type mysqlIngredient struct {
	db *gorm.DB
}

func NewMySQLIngredient(db *gorm.DB) *mysqlIngredient {
	return &mysqlIngredient{db}
}
