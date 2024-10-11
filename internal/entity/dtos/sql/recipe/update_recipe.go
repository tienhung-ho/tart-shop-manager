package recipemodel

import (
	"tart-shop-manager/internal/common"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

type UpdateRecipe struct {
	RecipeID          uint64                                   `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"-"`
	ProductID         uint64                                   `gorm:"column:product_id;not null" json:"product_id"`
	Size              string                                   `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost              float64                                  `gorm:"column:cost;not null" json:"cost"`
	Description       string                                   `gorm:"column:description;type:text" json:"description"`
	RecipeIngredients []recipeingredientmodel.RecipeIngredient `gorm:"foreignKey:RecipeID;references:RecipeID" json:"ingredients"`
	common.CommonFields
}

func (UpdateRecipe) TableName() string {
	return Recipe{}.TableName()
}
