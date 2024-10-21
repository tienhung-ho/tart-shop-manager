package supplyOrder

import (
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	"time"
)

type CreateSupplyOrderCache struct {
	SupplyOrderID uint64    `json:"supply_order_id"`
	OrderDate     time.Time `json:"order_date"`
	Description   string    `json:"description"`
	TotalAmount   float64   `json:"total_amount"`
	SupplierID    uint64    `json:"supplier_id"`
	common.CommonFields
	//Supplier      Supplier       `gorm:"foreignKey:SupplierID" json:"supplier"`
	SupplyOrderItems []supplyordermodel.SupplyOrderItem  `json:"supplyorder_item"`
	Ingredients      []supplyordermodel.CreateIngredient `json:"ingredients" binding:"required,dive,required" validate:"required"`
}

func ToCreateSupplyOrderCache(s supplyordermodel.SupplyOrder) *CreateSupplyOrderCache {
	return &CreateSupplyOrderCache{
		SupplyOrderID:    s.SupplyOrderID,
		OrderDate:        s.OrderDate,
		Description:      s.Description,
		TotalAmount:      s.TotalAmount,
		SupplierID:       s.SupplierID,
		CommonFields:     s.CommonFields,
		SupplyOrderItems: s.SupplyOrderItems,
		//Ingredients:      s.Ingredients,
	}
}
