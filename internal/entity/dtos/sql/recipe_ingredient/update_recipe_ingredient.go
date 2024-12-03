package recipeingredientmodel

type UpdateRecipeIngredient struct {
	RecipeID     uint64  `gorm:"column:recipe_id;primaryKey" json:"recipe_id"`
	IngredientID uint64  `gorm:"column:ingredient_id;primaryKey" json:"ingredient_id"`
	Unit         string  `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	Quantity     float64 `gorm:"column:quantity;not null" json:"quantity"`
}

func (UpdateRecipeIngredient) TableName() string {
	return "RecipeIngredient"
}
