package common

import (
	"gorm.io/gorm"
	"time"
)

type CommonFields struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	Status    *Status        `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status"`
	CreatedBy string         `gorm:"column:created_by;type:char(30)" json:"-"`
	UpdatedBy string         `gorm:"column:updated_by;type:char(30)" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
