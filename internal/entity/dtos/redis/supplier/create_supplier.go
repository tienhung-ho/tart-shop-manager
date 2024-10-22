package suppliercachemodel

import "tart-shop-manager/internal/common"

type CreateSupplier struct {
	SupplierID  uint64 `gorm:"primaryKey;autoIncrement;column:supplier_id" json:"-"`
	Name        string `gorm:"column:name;size:200;not null" json:"name"  validate:"required"`
	Description string `gorm:"column:description" json:"description"`
	ContactInfo string `gorm:"column:contactInfo;size:200;not null" json:"contact_info"  validate:"required"`
	Address     string `gorm:"column:address;size:200;not null" json:"address"  validate:"required"`
	common.CommonFields
}
