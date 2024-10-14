package supplyordermodel

import (
	"tart-shop-manager/internal/common"
)

var (
	EntityName = "SupplyOrder"
)

// SupplyOrder represents the supply order entity
type SupplyOrder struct {
	SupplyOrderID uint              `gorm:"column:supplyorder_id;primaryKey;autoIncrement" json:"supply_order_id"`
	OrderDate     common.CustomDate `gorm:"column:order_date;not null" json:"order_date"`
	Description   string            `gorm:"column:description;type:text" json:"description"`
	TotalAmount   float64           `gorm:"column:total_amount;not null" json:"total_amount"`
	SupplierID    uint              `gorm:"column:supplier_id;not null;index" json:"supplier_id"`
	common.CommonFields
	//Supplier      Supplier       `gorm:"foreignKey:SupplierID" json:"supplier"`
}

func (SupplyOrder) TableName() string {
	return "SupplyOrder"
}
