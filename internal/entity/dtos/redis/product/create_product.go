package productcachemodel

import (
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

type CreateProduct struct {
	ProductID       uint64                  `json:"product_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	QuantityInStock int                     `json:"quantity_in_stock"`
	ImageID         uint64                  `gorm:"column:image_id;size:300;foreignKey:ImageID;references:ImageID;not null" json:"image_id"`
	CategoryID      uint64                  `json:"category_id"`
	Category        *categorymodel.Category `json:"category"` // Liên kết với Category
	Recipes         []recipemodel.Recipe    `json:"recipes"`  // Một Product có nhiều Recipe
	*common.CommonFields
}
