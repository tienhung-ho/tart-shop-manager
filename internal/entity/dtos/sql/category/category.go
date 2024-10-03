package categorymodel

import (
	"tart-shop-manager/internal/common"
	categorycachemodel "tart-shop-manager/internal/entity/dtos/redis/category"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

var (
	EntityName   = "category"
	SelectFields = []string{
		"category_id",
		"name",
		"description",
		"status",
	}
)

var AllowedSortFields = map[string]bool{
	"name":        true,
	"created_at":  true,
	"updated_at":  true,
	"category_id": true,
}

type Category struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`
	Name        string             `gorm:"column:name;size:200;not null" json:"name"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `gorm:"foreignKey:CategoryID;references:CategoryID" json:"images"`
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
