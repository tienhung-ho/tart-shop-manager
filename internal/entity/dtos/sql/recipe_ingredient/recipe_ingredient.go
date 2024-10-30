package recipeingredientmodel

import ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"

type RecipeIngredient struct {
	RecipeID     uint64                     `gorm:"column:recipe_id;primaryKey" json:"recipe_id"`
	IngredientID uint64                     `gorm:"column:ingredient_id;primaryKey" json:"ingredient_id"`
	Quantity     float64                    `gorm:"column:quantity;not null" json:"quantity"`
	Ingredient   ingredientmodel.Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient"`
}

func (RecipeIngredient) TableName() string {
	return "RecipeIngredient"
}

func (r RecipeIngredient) ToUpdateRecipeIngredient() *UpdateRecipeIngredient {
	return &UpdateRecipeIngredient{
		RecipeID:     r.RecipeID,
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
	}
}
