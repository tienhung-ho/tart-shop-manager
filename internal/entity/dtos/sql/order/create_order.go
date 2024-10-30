package ordermodel

import (
	"tart-shop-manager/internal/common"
)

type CreateOrder struct {
	OrderID     uint64            `gorm:"column:order_id;primaryKey;autoIncrement:true" json:"order_id"`
	AccountID   uint64            `gorm:"column:account_id;not null" json:"account_id"`
	TotalAmount float64           `gorm:"column:total_amount;type:decimal(11,2);not null;default:0.00" json:"total_amount"`
	Tax         float64           `gorm:"column:tax;type:decimal(10,2);default:0.00" json:"tax"`
	OrderItems  []CreateOrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"order_items"`
	common.CommonFields
}

func (CreateOrder) TableName() string {
	return Order{}.TableName()
}
