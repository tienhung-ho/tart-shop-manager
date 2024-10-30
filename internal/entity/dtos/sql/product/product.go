package productmodel

import (
	"tart-shop-manager/internal/common"
	productcachemodel "tart-shop-manager/internal/entity/dtos/redis/product"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

var (
	EntityName   = "Product"
	SelectFields = []string{
		"product_id",
		"name",
		"description",
		"status",
		"quantity_in_stock",
	}
)

var AllowedSortFields = map[string]bool{
	"name":       true,
	"created_at": true,
	"updated_at": true,
	"product_id": true,
}

type Product struct {
	ProductID       uint64                  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	Name            string                  `gorm:"column:name;size:200;not null" json:"name"`
	Description     string                  `gorm:"column:description;type:text" json:"description"`
	QuantityInStock int                     `gorm:"column:quantity_in_stock;not null" json:"quantity_in_stock"`
	Price           float64                 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	Images          []imagemodel.Image      `gorm:"foreignKey:ProductID;references:ProductID" json:"images"`
	CategoryID      uint64                  `gorm:"column:category_id;not null" json:"category_id"`
	Category        *categorymodel.Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category"`
	Recipes         []recipemodel.Recipe    `gorm:"foreignKey:ProductID" json:"recipes"`
	*common.CommonFields
}

func (Product) TableName() string {
	return "Product"
}

func (p *Product) ToCreateAccount() *productcachemodel.CreateProduct {
	return &productcachemodel.CreateProduct{
		ProductID:       p.ProductID,
		Name:            p.Name,
		Description:     p.Description,
		QuantityInStock: p.QuantityInStock,
		CategoryID:      p.CategoryID,
		Recipes:         p.Recipes,
		Category:        p.Category,
		Images:          p.Images,
		CommonFields: &common.CommonFields{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			Status:    p.Status,
			CreatedBy: p.CreatedBy,
			UpdatedBy: p.UpdatedBy,
		},
	}
}
