package stockbatchcachemodel

import (
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	"time"
)

type CreateStockBatch struct {
	StockBatchID   uint64                      `gorm:"column:stockbatch_id;primaryKey;autoIncrement" json:"stockbatch_id"`
	Quantity       float64                     `gorm:"column:quantity;not null" json:"quantity" validate:"required,gt=0"`
	ExpirationDate time.Time                   `gorm:"column:expiration_date;not null" json:"expiration_date" validate:"required"`
	ReceivedDate   time.Time                   `gorm:"column:received_date;not null" json:"received_date" validate:"required"`
	IngredientID   uint64                      `gorm:"column:ingredient_id;not null" json:"ingredient_id" validate:"required"`
	Ingredient     *ingredientmodel.Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient" binding:"required,dive,required" validate:"required"`
	common.CommonFields
}
