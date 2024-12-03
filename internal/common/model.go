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
	Status    *Status        `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status,omitempty"`
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

		// Check if the destination object is a pointer and is valid
		if dest.Kind() == reflect.Ptr && dest.IsValid() {
			// Ensure it's a struct and it has the UpdatedBy field
			if dest.Elem().Kind() == reflect.Struct {
				if updatedByField := dest.Elem().FieldByName("UpdatedBy"); updatedByField.IsValid() && updatedByField.CanSet() {
					// Set the value of UpdatedBy
					updatedByField.SetString(email)
					// Ensure the update is persisted
					tx.Statement.SetColumn("updated_by", email)
				}
			} else {
				log.Print("Destination is not a struct")
			}
		} else {
			log.Print("Destination is not a pointer or is invalid")
			//return errors.New("invalid destination object")
		}
	} else {
		log.Print("Email is missing from context")
	}

	return nil
}
