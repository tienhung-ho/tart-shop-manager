package suppliermodel

import (
	"tart-shop-manager/internal/common"
	suppliercachemodel "tart-shop-manager/internal/entity/dtos/redis/supplier"
)

var (
	EntityName = "Supplier"
)

type Supplier struct {
	SupplierID  uint64 `gorm:"primaryKey;autoIncrement;column:supplier_id" json:"supplier_id"`
	Name        string `gorm:"column:name;size:200;not null" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	ContactInfo string `gorm:"column:contactInfo;size:200;not null" json:"contact_info"`
	Address     string `gorm:"column:address;size:200;not null" json:"address"`
	common.CommonFields
}

func (Supplier) TableName() string {
	return "Supplier"
}

func (s Supplier) ToCreateSupplierCache() *suppliercachemodel.CreateSupplier {
	return &suppliercachemodel.CreateSupplier{
		SupplierID:   s.SupplierID,
		Name:         s.Name,
		Description:  s.Description,
		ContactInfo:  s.ContactInfo,
		Address:      s.Address,
		CommonFields: s.CommonFields,
	}
}
