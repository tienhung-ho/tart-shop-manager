package reportmodel

import (
	"time"
)

var (
	SupplyOrderReportName = "SupplyOrderReport"
)

type StockBatch struct {
	StockBatchID   uint64    `gorm:"column:stockbatch_id;primaryKey;autoIncrement" json:"stockbatch_id"`
	Quantity       float64   `gorm:"column:quantity;not null" json:"quantity"`
	ExpirationDate time.Time `gorm:"column:expiration_date;not null" json:"expiration_date"`
	ReceivedDate   time.Time `gorm:"column:received_date;not null" json:"received_date"`
	IngredientID   uint64    `gorm:"column:ingredient_id;not null" json:"ingredient_id"`
}

func (StockBatch) TableName() string {
	return "StockBatch"
}

type Ingredient struct {
	IngredientID uint64 `gorm:"column:ingredient_id;primaryKey;autoIncrement" json:"-"`
	Name         string `gorm:"column:name;size:200;not null" json:"name"`
}

func (Ingredient) TableName() string {
	return "Ingredient"
}

type SupplyOrderItem struct {
	SupplyOrderItemID uint64     `gorm:"column:supplyorderitem_id;primaryKey;autoIncrement" json:"supplyorderitem_id"`
	Price             float64    `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity          float64    `gorm:"column:quantity;not null" json:"quantity"`
	Unit              string     `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	SupplyOrderID     uint64     `gorm:"column:supplyorder_id;unique;not null;index" json:"-"`
	IngredientID      uint64     `gorm:"column:ingredient_id;" json:"ingredient_id"`
	Ingredient        Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient"`
}

func (SupplyOrderItem) TableName() string {
	return "SupplyOrderItem"
}

type SupplyOrder struct {
	SupplyOrderID    uint64            `gorm:"column:supplyorder_id;primaryKey;autoIncrement" json:"supply_order_id"`
	OrderDate        time.Time         `gorm:"column:order_date;not null" json:"order_date"`
	TotalAmount      float64           `gorm:"column:total_amount;not null" json:"total_amount"`
	SupplierID       uint64            `gorm:"column:supplier_id;not null;index" json:"supplier_id"`
	Supplier         Supplier          `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier"`
	SupplyOrderItems []SupplyOrderItem `gorm:"foreignKey:SupplyOrderID;references:SupplyOrderID" json:"supplierorder_item"`
}

func (SupplyOrder) TableName() string {
	return "SupplyOrder"
}

type Supplier struct {
	SupplierID  uint64 `gorm:"primaryKey;autoIncrement;column:supplier_id" json:"-"`
	Name        string `gorm:"column:name;size:200;not null" json:"name"`
	ContactInfo string `gorm:"column:contactInfo;size:200;not null" json:"contact_info"`
}

func (Supplier) TableName() string {
	return "Supplier"
}
