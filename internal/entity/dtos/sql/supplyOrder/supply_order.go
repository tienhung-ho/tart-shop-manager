package supplyordermodel

import (
	"tart-shop-manager/internal/common"
	"time"
)

var (
	EntityName     = "SupplyOrder"
	ItemEntityName = "SupplyOrderItem"
)

// SupplyOrder represents the supply order entity
type SupplyOrder struct {
	SupplyOrderID uint64    `gorm:"column:supplyorder_id;primaryKey;autoIncrement" json:"supply_order_id"`
	OrderDate     time.Time `gorm:"column:order_date;not null" json:"order_date"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
	TotalAmount   float64   `gorm:"column:total_amount;not null" json:"total_amount"`
	SupplierID    uint64    `gorm:"column:supplier_id;not null;index" json:"supplier_id"`
	common.CommonFields
	SupplyOrderItems []SupplyOrderItem `gorm:"foreignKey:SupplyOrderID;references:SupplyOrderID" json:"supplierorder_item"`
	//Supplier      Supplier       `gorm:"foreignKey:SupplierID" json:"supplier"`
	//Ingredients []CreateIngredient `gorm:"-" json:"ingredients" binding:"required,dive,required" validate:"required"`
}

func (SupplyOrder) TableName() string {
	return "SupplyOrder"
}
