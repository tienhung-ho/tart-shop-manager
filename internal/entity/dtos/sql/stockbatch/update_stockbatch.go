package stockbatchmodel

import "tart-shop-manager/internal/common"

type UpdateStockBatch struct {
	StockBatchID   uint64            `gorm:"column:stockbatch_id;primaryKey;autoIncrement" json:"-"`
	Quantity       int               `gorm:"column:quantity;not null" json:"quantity" validate:"required,gt=0"`
	ExpirationDate common.CustomDate `gorm:"column:expiration_date;not null" json:"expiration_date" validate:"required"`
	ReceivedDate   common.CustomDate `gorm:"column:received_date;not null" json:"received_date" validate:"required"`
	IngredientID   uint64            `gorm:"column:ingredient_id;not null" json:"ingredient_id" validate:"required"`
	common.CommonFields
}

func (UpdateStockBatch) TableName() string {
	return StockBatch{}.TableName()
}
