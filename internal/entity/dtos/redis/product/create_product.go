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
	ImageURL        string                  `json:"image_url"`
	CategoryID      uint64                  `json:"category_id"`
	Category        *categorymodel.Category `json:"category"` // Liên kết với Category
	Recipes         []recipemodel.Recipe    `json:"recipes"`  // Một Product có nhiều Recipe
	*common.CommonFields
}
