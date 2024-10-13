package common

import (
	"gorm.io/gorm"
	"log"
	"reflect"
	"time"
)

type CommonFields struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	Status    *Status        `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status"`
	CreatedBy string         `gorm:"column:created_by;type:char(30);default:'system'" json:"-"`
	UpdatedBy string         `gorm:"column:updated_by;type:char(30)" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Hook BeforeCreate để thiết lập CreatedBy từ context
func (cf *CommonFields) BeforeCreate(tx *gorm.DB) (err error) {
	if email, ok := tx.Statement.Context.Value("email").(string); ok {
		cf.CreatedBy = email
	} else {
		//log.Printf("Email is missing from context")
	}
	return nil
}

// Hook BeforeUpdate để thiết lập UpdatedBy từ context
func (cf *CommonFields) BeforeUpdate(tx *gorm.DB) (err error) {
	if email, ok := tx.Statement.Context.Value("email").(string); ok {
		cf.UpdatedBy = email

		log.Print("Updating UpdatedBy field")
		// Get the destination object and handle it with reflect.Indirect to avoid pointer issues
		dest := reflect.Indirect(reflect.ValueOf(tx.Statement.Dest))

		// Check if the struct has an UpdatedBy field
		if dest.Kind() == reflect.Struct {
			if updatedByField := dest.FieldByName("UpdatedBy"); updatedByField.IsValid() && updatedByField.CanSet() {
				// Set the value
				updatedByField.SetString(email)

				// Use GORM's Update method to ensure the change is persisted
				tx.Statement.SetColumn("updated_by", email)
			}
		} else {
			log.Print("The destination object is not a struct")
		}
	} else {
		log.Print("Email is missing from context")
	}

	return nil
}
