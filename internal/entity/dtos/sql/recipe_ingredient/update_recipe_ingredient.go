package recipeingredientmodel

type UpdateRecipeIngredient struct {
	RecipeID     uint64  `gorm:"column:recipe_id;primaryKey" json:"recipe_id"`
	IngredientID uint64  `gorm:"column:ingredient_id;primaryKey" json:"ingredient_id"`
	Quantity     float64 `gorm:"column:quantity;not null" json:"quantity"`
}

func (UpdateRecipeIngredient) TableName() string {
	return "RecipeIngredient"
}
