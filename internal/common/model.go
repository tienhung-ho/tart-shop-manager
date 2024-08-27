package common

import "time"

type CommonFields struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Status    *Status   `gorm:"column:status;type:enum('Pending', 'Active', 'Inactive');default:Pending" json:"status"`
	CreatedBy string    `gorm:"column:created_by;type:char(30)" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by;type:char(30)" json:"updated_by"`
}
