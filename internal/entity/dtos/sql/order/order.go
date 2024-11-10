package ordermodel

import (
	"tart-shop-manager/internal/common"
	ordercachemodel "tart-shop-manager/internal/entity/dtos/redis/order"
	orderitemmodel "tart-shop-manager/internal/entity/dtos/sql/orderItem"
)

var (
	EntityName = "Order"
)

type Product struct {
}

type Order struct {
	OrderID     uint64            `gorm:"column:order_id;primaryKey;autoIncrement:true" json:"order_id"`
	AccountID   uint64            `gorm:"column:account_id;not null" json:"account_id"`
	TotalAmount float64           `gorm:"column:total_amount;type:decimal(11,2);not null;default:0.00" json:"total_amount"`
	Tax         float64           `gorm:"column:tax;type:decimal(10,2);default:0.00" json:"tax"`
	OrderDate   common.CustomDate `gorm:"column:order_date;not null" json:"order_date"`
	//Recipes     []recipemodel.Recipe `gorm:"many2many:OrderRecipe;foreignKey:OrderID;joinForeignKey:OrderID;References:RecipeID;joinReferences:RecipeID"`
	OrderItems []orderitemmodel.OrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"order_items,omitempty"`
	common.CommonFields
}

type RecipeIngredientQuantity struct {
	IngredientID int64
	Quantity     int64
}

func (Order) TableName() string {
	return "Order"
}

func (o *Order) ToCreateOrder() *ordercachemodel.CreateOrder {
	return &ordercachemodel.CreateOrder{
		OrderID:      o.OrderID,
		AccountID:    o.AccountID,
		TotalAmount:  o.TotalAmount,
		Tax:          o.Tax,
		OrderItems:   o.OrderItems,
		CommonFields: o.CommonFields,
		OrderDate:    o.OrderDate,
	}
}
