package productmodel

import (
	"tart-shop-manager/internal/common"
)

type CreateProduct struct {
	ProductID       uint64 `gorm:"column:product_id;primaryKey;autoIncrement" json:"-"`
	Name            string `gorm:"column:name;size:200;not null" json:"name" validate:"required"`
	Description     string `gorm:"column:description;type:text" json:"description"`
	QuantityInStock int    `gorm:"column:quantity_in_stock;not null" json:"quantity_in_stock" validate:"required,gt=0"`
	ImageURL        string `gorm:"column:image_url;size:300;not null" json:"image_url"`
	CategoryID      uint64 `gorm:"column:category_id;not null" json:"category_id" validate:"required"`
	//Category        *categorymodel.Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category"`
	//Recipes         []recipemodel.Recipe    `gorm:"foreignKey:ProductID" json:"recipes"`
	*common.CommonFields
}

func (CreateProduct) TableName() string {
	return "Product"
}
