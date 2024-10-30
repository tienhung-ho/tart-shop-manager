package supplyordermodel

import "tart-shop-manager/internal/common"

type UpdateSupplyOrderItem struct {
	SupplyOrderItemID uint64  `gorm:"column:supplyorderitem_id;primaryKey;autoIncrement" json:"supplyorderitem_id"`
	Price             float64 `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity          float64 `gorm:"column:quantity;not null" json:"quantity"`
	Unit              string  `gorm:"column:unit;type:varchar(200);not null" json:"unit"`
	IngredientID      uint64  `gorm:"column:ingredient_id;not null;index" json:"ingredient_id"`
	SupplyOrderID     uint64  `gorm:"column:supplyorder_id;unique;not null;index" json:"supplyorder_id"`
	StockBatchID      uint64  `gorm:"column:stockbatch_id;not null;index" json:"stockbatch_id"`
	common.CommonFields
}

func (UpdateSupplyOrderItem) TableName() string {
	return SupplyOrderItem{}.TableName()
}
