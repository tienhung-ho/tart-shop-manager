package ingredientcachemodel

import "tart-shop-manager/internal/common"

type CreateIngredient struct {
	IngredientID uint   `gorm:"column:ingredient_id;primaryKey;autoIncrement" json:"ingredient_id"`
	Name         string `gorm:"column:name;size:200;not null" json:"name"`
	Description  string `gorm:"column:description;type:text" json:"description"`
	Unit         string `gorm:"column:unit;size:100;not null" json:"unit"`
	common.CommonFields
}
