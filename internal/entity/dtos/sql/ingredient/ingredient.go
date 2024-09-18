package ingredientmodel

import (
	"tart-shop-manager/internal/common"
	ingredientcachemodel "tart-shop-manager/internal/entity/dtos/redis/ingredient"
)

var (
	EntityName = "ingredient"
)

type Ingredient struct {
	IngredientID uint   `gorm:"column:ingredient_id;primaryKey;autoIncrement" json:"ingredient_id"`
	Name         string `gorm:"column:name;size:200;not null" json:"name"`
	Description  string `gorm:"column:description;type:text" json:"description"`
	Unit         string `gorm:"column:unit;size:100;not null" json:"unit"`
	common.CommonFields
}

func (Ingredient) TableName() string {
	return "Ingredient"
}

func (i Ingredient) ToCreateIngredientCache() *ingredientcachemodel.CreateIngredient {
	return &ingredientcachemodel.CreateIngredient{
		IngredientID: i.IngredientID,
		Name:         i.Name,
		Description:  i.Description,
		Unit:         i.Unit,
		CommonFields: i.CommonFields,
	}
}
