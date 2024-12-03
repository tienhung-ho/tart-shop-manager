package recipemodel

import (
	"tart-shop-manager/internal/common"
	recipecachemodel "tart-shop-manager/internal/entity/dtos/redis/recipe"
	recipeingredientmodel "tart-shop-manager/internal/entity/dtos/sql/recipe_ingredient"
)

var (
	EntityName = "Recipe"
)

type product struct {
	ProductID   uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	Name        string  `gorm:"column:name;size:200;not null" json:"name"`
	Price       float64 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	Description string  `gorm:"column:description;type:text" json:"description"`
	ImageURL    string  `gorm:"column:image_url;size:300;not null" json:"image_url"`
	CategoryID  uint64  `gorm:"column:category_id;not null" json:"category_id"`
	*common.CommonFields
}

func (product) TableName() string {
	return "Product"
}

type Recipe struct {
	RecipeID          uint64                                   `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	ProductID         uint64                                   `gorm:"column:product_id;not null" json:"product_id"`
	Product           *product                                 `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"` // Liên kết với Product
	Size              string                                   `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost              float64                                  `gorm:"column:cost;not null" json:"cost"`
	Description       string                                   `gorm:"column:description;type:text" json:"description"`
	RecipeIngredients []recipeingredientmodel.RecipeIngredient `gorm:"foreignKey:RecipeID;references:RecipeID" json:"ingredients,omitempty"`
	common.CommonFields
}

func (Recipe) TableName() string {
	return "Recipe"
}

func (r Recipe) ToCreateRecipe() *recipecachemodel.CreateRecipe {
	return &recipecachemodel.CreateRecipe{
		RecipeID:  r.RecipeID,
		ProductID: r.ProductID,
		Product: &recipecachemodel.RecipeProduct{
			ProductID:    r.ProductID,
			Name:         r.Product.Name,
			Price:        r.Product.Price,
			Description:  r.Product.Description,
			ImageURL:     r.Product.ImageURL,
			CategoryID:   r.Product.CategoryID,
			CommonFields: r.Product.CommonFields,
		},
		Size:              r.Size,
		Cost:              r.Cost,
		Description:       r.Description,
		RecipeIngredients: r.RecipeIngredients,
		CommonFields:      r.CommonFields,
	}
}
