package categorymodel

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type CreateCategory struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"-"`
	Name        string             `gorm:"column:name;size:200;not null;unique" json:"name" validate:"required"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `gorm:"foreignKey:CategoryID;references:CategoryID" json:"images"`
	common.CommonFields
}

func (CreateCategory) TableName() string {
	return Category{}.TableName()
}
