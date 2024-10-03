package common

import (
	"gorm.io/gorm"
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

		// Get the destination object
		dest := reflect.ValueOf(tx.Statement.Dest).Elem()

		// Check if the struct has an UpdatedBy field
		if updatedByField := dest.FieldByName("UpdatedBy"); updatedByField.IsValid() && updatedByField.CanSet() {
			// Set the value
			updatedByField.SetString(email)

			// Use GORM's Update method to ensure the change is persisted
			tx.Update("updated_by", email)
		}
	} else {
		// Consider adding proper error handling or logging here
		// log.Printf("Email is missing from context")
	}

	return nil
}
