package recipeingredientmodel

type RecipeIngredientCreate struct {
	RecipeID     uint64  `gorm:"column:recipe_id;primaryKey" json:"recipe_id"`
	IngredientID uint64  `gorm:"column:ingredient_id;primaryKey" json:"ingredient_id"`
	Quantity     float64 `gorm:"column:quantity;not null" json:"quantity"`
}

func (RecipeIngredientCreate) TableName() string {
	return RecipeIngredient{}.TableName()
}
