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
	Status      string `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status"`
	//Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // Một Category có nhiều Product
	*common.CommonFields
}

func (Category) TableName() string {
	return "Category"
}

func (c Category) ToCreateCategoryCache() *categorycachemodel.CreateCategory {
	return &categorycachemodel.CreateCategory{
		CategoryID: c.CategoryID,
	}
}
