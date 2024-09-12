package categorymodel

import "tart-shop-manager/internal/common"

type CreateCategory struct {
	CategoryID  uint64 `gorm:"column:category_id;primaryKey;autoIncrement" json:"-"`
	Name        string `gorm:"column:name;size:200;not null;unique" json:"name" validate:"required"`
	Description string `gorm:"column:description;type:text" json:"description"`
	//Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // Một Category có nhiều Product
	*common.CommonFields
}

func (CreateCategory) TableName() string {
	return Category{}.TableName()
}
