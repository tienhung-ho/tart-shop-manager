package categorycachemodel

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type CreateCategory struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`
	Name        string             `gorm:"column:name;size:200;not null" json:"name"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `json:"images"`
	//Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // Một Category có nhiều Product
	*common.CommonFields
}

func (CreateCategory) TableName() string {
	return "Category"
}
