package productmodel

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type CreateProduct struct {
	ProductID   uint64 `gorm:"column:product_id;primaryKey;autoIncrement" json:"-"`
	Name        string `gorm:"column:name;size:200;not null" json:"name" validate:"required"`
	Description string `gorm:"column:description;type:text" json:"description"`
	//QuantityInStock     int                      `gorm:"column:quantity_in_stock;not null" json:"quantity_in_stock" validate:"required,gt=0"`
	Images              []imagemodel.UpdateImage `gorm:"-" json:"images"`
	Price               float64                  `gorm:"column:price;type:decimal(11,2)" json:"price"`
	CategoryID          uint64                   `gorm:"column:category_id;not null" json:"category_id" validate:"required"`
	common.CommonFields `gorm:"embedded"`
}

func (CreateProduct) TableName() string {
	return Product{}.TableName()
}
