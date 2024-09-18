package stockbatchstorage

import "gorm.io/gorm"

type mysqlStockBatch struct {
	db *gorm.DB
}

func NewMySQLStockBatch(db *gorm.DB) *mysqlStockBatch {
	return &mysqlStockBatch{db}
}
