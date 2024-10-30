package productcachemodel

import (
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

type CreateProduct struct {
	ProductID       uint64                  `json:"product_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	QuantityInStock int                     `json:"quantity_in_stock"`
	Price           float64                 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	Images          []imagemodel.Image      `gorm:"foreignKey:ProductID;references:ProductID" json:"images"`
	CategoryID      uint64                  `json:"category_id"`
	Category        *categorymodel.Category `json:"category"` // Liên kết với Category
	Recipes         []recipemodel.Recipe    `json:"recipes"`  // Một Product có nhiều Recipe
	*common.CommonFields
}
