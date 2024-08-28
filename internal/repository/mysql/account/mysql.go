package accountstorage

import "gorm.io/gorm"

type mysqlAccount struct {
	db *gorm.DB
}

func NewMySQLAccount(db *gorm.DB) *mysqlAccount {
	return &mysqlAccount{db}
}
