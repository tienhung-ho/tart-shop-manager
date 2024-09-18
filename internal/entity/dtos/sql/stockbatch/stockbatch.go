package stockbatchmodel

import (
	"tart-shop-manager/internal/common"
	"time"
)

type StockBatch struct {
	StockBatchID   uint      `gorm:"column:stockbatch_id;primaryKey;autoIncrement" json:"stockbatch_id"`
	Quantity       int       `gorm:"column:quantity;not null" json:"quantity"`
	ExpirationDate time.Time `gorm:"column:expiration_date;not null" json:"expiration_date"`
	ReceivedDate   time.Time `gorm:"column:received_date;not null" json:"received_date"`
	IngredientID   uint      `gorm:"column:ingredient_id;not null" json:"ingredient_id"`
	common.CommonFields
}

func (StockBatch) TableName() string {
	return "StockBatch"
}
