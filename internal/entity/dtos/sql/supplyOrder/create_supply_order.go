package supplyordermodel

import (
	"tart-shop-manager/internal/common"
)

type CreateIngredient struct {
	IngredientID   uint64            `json:"ingredient_id"`
	Unit           string            `json:"unit"`
	Quantity       int               `json:"quantity"`
	Price          float64           `json:"price"`
	ExpirationDate common.CustomDate `json:"expiration_date"`
	ReceivedDate   common.CustomDate `json:"received_date"`
}

type CreateSupplyOrder struct {
	SupplyOrderID uint               `gorm:"column:supplyorder_id;primaryKey;autoIncrement" json:"supply_order_id"`
	OrderDate     common.CustomDate  `gorm:"column:order_date;not null" json:"order_date" validate:"required"`
	Description   string             `gorm:"column:description;type:text" json:"description" validate:"required"`
	TotalAmount   float64            `gorm:"column:total_amount;not null" json:"total_amount" `
	SupplierID    uint               `gorm:"column:supplier_id;not null;index" json:"supplier_id" validate:"required"`
	Ingredients   []CreateIngredient `gorm:"-" json:"ingredients" binding:"required,dive,required" validate:"required"`
	common.CommonFields
	//SupplierOrderItem SupplyOrderItem `gorm:"foreignKey:SupplyOrderItemID" json:"supplier"`
}

func (CreateSupplyOrder) TableName() string {
	return SupplyOrder{}.TableName()
}
