package stockbatchmodel

import (
	"tart-shop-manager/internal/common"
	stockbatchcachemodel "tart-shop-manager/internal/entity/dtos/redis/stockBatch"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	"time"
)

var (
	EntityName = "StockBatch"
)

type StockBatch struct {
	StockBatchID   uint64                      `gorm:"column:stockbatch_id;primaryKey;autoIncrement" json:"stockbatch_id"`
	Quantity       int                         `gorm:"column:quantity;not null" json:"quantity"`
	ExpirationDate time.Time                   `gorm:"column:expiration_date;not null" json:"expiration_date"`
	ReceivedDate   time.Time                   `gorm:"column:received_date;not null" json:"received_date"`
	IngredientID   uint64                      `gorm:"column:ingredient_id;not null" json:"ingredient_id"`
	Ingredient     *ingredientmodel.Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient" binding:"required,dive,required" validate:"required"`
	common.CommonFields
}

func (StockBatch) TableName() string {
	return "StockBatch"
}

func (s StockBatch) ToCreateStockBatch() *stockbatchcachemodel.CreateStockBatch {
	return &stockbatchcachemodel.CreateStockBatch{
		StockBatchID:   s.StockBatchID,
		Quantity:       s.Quantity,
		ExpirationDate: s.ExpirationDate,
		ReceivedDate:   s.ReceivedDate,
		IngredientID:   s.IngredientID,
		CommonFields:   s.CommonFields,
	}
}
