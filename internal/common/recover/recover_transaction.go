package commonrecover

import (
	"fmt"
	"gorm.io/gorm"
)

func RecoverTransaction(db *gorm.DB) {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from panic: %v\n", r)
		db.Rollback()
		panic(r) // Optionally re-panic to handle further up the stack
	}
}
