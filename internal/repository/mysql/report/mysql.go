package reportstorage

import "gorm.io/gorm"

var (
	SelectFields = []string{
		"recipe_id",
		"product_id",
		"size",
		"cost",
		"description",
		"status",
	}
	AllowedSortFields = map[string]bool{
		"recipe_id":  true,
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"cost":       true,
	}
)

type mysqlReportOrder struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlReportOrder {
	return &mysqlReportOrder{db}
}
