package recipemodel

import "tart-shop-manager/internal/common"

type CreateRecipe struct {
	RecipeID    uint64  `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"-"`
	ProductID   uint64  `gorm:"column:product_id;not null" json:"product_id"`
	Size        string  `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost        float64 `gorm:"column:cost;not null" json:"cost"`
	Description string  `gorm:"column:description;type:text" json:"description"`
	common.CommonFields
}

func (CreateRecipe) TableName() string {
	return Recipe{}.TableName()
}
