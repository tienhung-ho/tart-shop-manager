package stockbatchstorage

import "gorm.io/gorm"

type mysqlStockBatch struct {
	db *gorm.DB
}

var (
	AllowedSortFields = map[string]bool{
		"stockbatch_id":   true,
		"quantity":        true,
		"expiration_date": true,
		"received_date":   true,
		"updated_at":      true,
	}
)

func NewMySQLStockBatch(db *gorm.DB) *mysqlStockBatch {
	return &mysqlStockBatch{db}
}
