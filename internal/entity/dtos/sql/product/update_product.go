package productmodel

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type UpdateProduct struct {
	ProductID   uint64             `gorm:"column:product_id;primaryKey;autoIncrement" json:"-"`
	Name        string             `gorm:"column:name;size:200;not null" json:"name" validate:"required"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `gorm:"foreignKey:ProductID;references:ProductID" json:"images"`
	CategoryID  uint64             `gorm:"column:category_id;not null" json:"category_id" validate:"required"`
	Price       float64            `gorm:"column:price;type:decimal(11,2)" json:"price"`
	common.CommonFields
}

func (UpdateProduct) TableName() string {
	return Product{}.TableName()
}
