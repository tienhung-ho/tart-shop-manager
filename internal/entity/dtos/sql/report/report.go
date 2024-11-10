package reportmodel

import (
	"tart-shop-manager/internal/common"
)

type RevenueReport struct {
	StartDate    *common.CustomDate `json:"startDate"`
	EndDate      *common.CustomDate `json:"endDate"`
	TotalRevenue float64            `json:"totalRevenue"`
	TotalCost    float64            `json:"totalCost"`
	TotalOrders  int                `json:"totalOrders"`
	Orders       []OrderSummary     `json:"orders"`
}

type OrderSummary struct {
	OrderID     uint64             `json:"orderId"`
	OrderDate   *common.CustomDate `json:"orderDate"`
	TotalAmount float64            `json:"totalAmount"`
	Items       []OrderItemSummary `json:"items"`
}

type OrderItemSummary struct {
	ProductID uint64  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	RecipeID  uint64  `json:"recipe_id"`
	Size      string  `json:"size"`
	Cost      float64 `json:"cost"`
}

type SupplyReport struct {
	StartDate         *common.CustomDate   `json:"startDate"`
	EndDate           *common.CustomDate   `json:"endDate"`
	TotalSupplyCost   float64              `json:"totalSupplyCost"`
	TotalSupplyOrders int                  `json:"totalSupplyOrders"`
	SupplyOrders      []SupplyOrderSummary `json:"supplyOrders"`
}

type SupplyOrderSummary struct {
	SupplyOrderID uint64                   `json:"supplyOrderId"`
	OrderDate     *common.CustomDate       `json:"orderDate"`
	SupplierID    uint64                   `json:"supplierId"`
	SupplierName  string                   `json:"supplierName"`
	TotalAmount   float64                  `json:"totalAmount"`
	ContactInfo   string                   `gorm:"column:contactInfo;size:200;not null" json:"contact_info"`
	Items         []SupplyOrderItemSummary `json:"items"`
}

type SupplyOrderItemSummary struct {
	IngredientID   uint64  `json:"ingredientId"`
	IngredientName string  `json:"ingredient_name"`
	Price          float64 `json:"price"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	TotalCost      float64 `json:"totalCost"`
}
