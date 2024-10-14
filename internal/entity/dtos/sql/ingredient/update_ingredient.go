package ingredientmodel

import "tart-shop-manager/internal/common"

type UpdateIngredient struct {
	IngredientID uint   `gorm:"column:ingredient_id;primaryKey;autoIncrement" json:"-"`
	Name         string `gorm:"column:name;size:200;not null" json:"name"`
	Description  string `gorm:"column:description;type:text" json:"description"`
	Unit         string `gorm:"column:unit;size:100;not null" json:"unit"`
	common.CommonFields
}

func (UpdateIngredient) TableName() string {
	return Ingredient{}.TableName()
}