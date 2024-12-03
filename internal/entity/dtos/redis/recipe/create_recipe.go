package recipecachemodel

import (
	"tart-shop-manager/internal/common"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

type RecipeProduct struct {
	ProductID   uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	Name        string  `gorm:"column:name;size:200;not null" json:"name"`
	Price       float64 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	Description string  `gorm:"column:description;type:text" json:"description"`
	ImageURL    string  `gorm:"column:image_url;size:300;not null" json:"image_url"`
	CategoryID  uint64  `gorm:"column:category_id;not null" json:"category_id"`
	*common.CommonFields
}

type CreateRecipe struct {
	RecipeID          uint64                                   `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	ProductID         uint64                                   `gorm:"column:product_id;not null" json:"product_id"`
	Unit              string                                   `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	Size              string                                   `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost              float64                                  `gorm:"column:cost;not null" json:"cost"`
	Product           *RecipeProduct                           `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"` // Liên kết với Product
	Description       string                                   `gorm:"column:description;type:text" json:"description"`
	RecipeIngredients []recipeingredientmodel.RecipeIngredient `gorm:"-" json:"ingredients"`
	common.CommonFields
}
