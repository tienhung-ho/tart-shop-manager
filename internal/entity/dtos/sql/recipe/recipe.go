package recipemodel

import (
	"tart-shop-manager/internal/common"
	recipecachemodel "tart-shop-manager/internal/entity/dtos/redis/recipe"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

var (
	EntityName   = "Recipe"
	SelectFields = []string{
		"recipe_id",
		"product_id",
		"size",
		"cost",
		"description",
		"status",
	}
)

type product struct {
	ProductID       uint64                  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	Name            string                  `gorm:"column:name;size:200;not null" json:"name"`
	Description     string                  `gorm:"column:description;type:text" json:"description"`
	QuantityInStock int                     `gorm:"column:quantity_in_stock;not null" json:"quantity_in_stock"`
	ImageURL        string                  `gorm:"column:image_url;size:300;not null" json:"image_url"`
	CategoryID      uint64                  `gorm:"column:category_id;not null" json:"category_id"`
	Category        *categorymodel.Category `gorm:"foreignKey:CategoryID" json:"category"` // Liên kết với Category
	*common.CommonFields
}

type Recipe struct {
	RecipeID    uint64   `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	ProductID   uint64   `gorm:"column:product_id;not null" json:"product_id"`
	Product     *product `gorm:"foreignKey:ProductID" json:"product"` // Liên kết với Product
	Size        string   `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost        float64  `gorm:"column:cost;not null" json:"cost"`
	Description string   `gorm:"column:description;type:text" json:"description"`
	common.CommonFields
}

func (Recipe) TableName() string {
	return "Recipe"
}

func (r Recipe) ToCreateRecipe() *recipecachemodel.CreateRecipe {
	return &recipecachemodel.CreateRecipe{
		RecipeID:     r.RecipeID,
		ProductID:    r.ProductID,
		Size:         r.Size,
		Cost:         r.Cost,
		Description:  r.Description,
		CommonFields: r.CommonFields,
	}
}
