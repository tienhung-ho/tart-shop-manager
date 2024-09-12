package categorymodel

import (
	"tart-shop-manager/internal/common"
	categorycachemodel "tart-shop-manager/internal/entity/dtos/redis/category"
)

var (
	EntityName = "category"
)

type Category struct {
	CategoryID  uint64 `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`
	Name        string `gorm:"column:name;size:200;not null" json:"name"`
	Description string `gorm:"column:description;type:text" json:"description"`
	//Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // Một Category có nhiều Product
	*common.CommonFields
}

func (Category) TableName() string {
	return "Category"
}

func (c Category) ToCreateCategoryCache() *categorycachemodel.CreateCategory {
	return &categorycachemodel.CreateCategory{
		CategoryID:  c.CategoryID,
		Name:        c.Name,
		Description: c.Description,
		CommonFields: &common.CommonFields{
			Status: c.Status,
		},
	}
}
