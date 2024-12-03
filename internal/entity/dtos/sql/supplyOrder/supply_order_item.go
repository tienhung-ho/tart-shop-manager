package supplyordermodel

import (
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

type CreateIngredient struct {
	IngredientID   uint64            `json:"ingredient_id"`
	Unit           string            `json:"unit"`
	Quantity       float64           `json:"quantity"`
	Price          float64           `json:"price"`
	ExpirationDate common.CustomDate `json:"expiration_date"`
	ReceivedDate   common.CustomDate `json:"received_date"`
}

type SupplyOrderItem struct {
	SupplyOrderItemID uint64  `gorm:"column:supplyorderitem_id;primaryKey;autoIncrement" json:"supplyorderitem_id"`
	Price             float64 `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity          int     `gorm:"column:quantity;not null" json:"quantity"`
	Unit              string  `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	IngredientID      uint64  `gorm:"column:ingredient_id;not null;index" json:"ingredient_id"`
	SupplyOrderID     uint64  `gorm:"column:supplyorder_id;unique;not null;index" json:"supplyorder_id"`
	StockBatchID      uint64  `gorm:"column:stockbatch_id;not null;index" json:"stockbatch_id"`

	// Quan hệ với các bảng khác
	//SupplyOrder SupplyOrder                `gorm:"belongsTo:SupplyOrder;foreignKey:SupplyOrderID;references:SupplyOrderID" json:"supply_order"`
	StockBatch *stockbatchmodel.StockBatch `gorm:"foreignKey:StockBatchID;references:StockBatchID" json:"stock_batch"`

	// Trường chung
	common.CommonFields
}

type SimpleSupplyOrderItemIngredient struct {
	Price        float64 `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Unit         string  `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	IngredientID uint64  `gorm:"column:ingredient_id;not null;index" json:"ingredient_id"`
}

func (SimpleSupplyOrderItemIngredient) TableName() string {
	return "SupplyOrderItem"
}

func (SupplyOrderItem) TableName() string {
	return "SupplyOrderItem"
}
