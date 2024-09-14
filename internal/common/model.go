package common

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type CommonFields struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	Status    *Status        `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status"`
	CreatedBy string         `gorm:"column:created_by;type:char(30);default:'system'" json:"-"`
	UpdatedBy string         `gorm:"column:updated_by;type:char(30);default:'system'" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Hook BeforeCreate để thiết lập CreatedBy từ context
func (cf *CommonFields) BeforeCreate(tx *gorm.DB) (err error) {
	if email, ok := tx.Statement.Context.Value("email").(string); ok {
		cf.CreatedBy = email
	} else {
		log.Printf("Email is missing from context")
	}
	return nil
}

// Hook BeforeUpdate để thiết lập UpdatedBy từ context
func (cf *CommonFields) BeforeUpdate(tx *gorm.DB) (err error) {
	if email, ok := tx.Statement.Context.Value("email").(string); ok {
		cf.UpdatedBy = email
	} else {
		log.Printf("Email is missing from context")
	}

	return nil
}
