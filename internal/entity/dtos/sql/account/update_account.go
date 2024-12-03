package accountmodel

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type UpdateAccount struct {
	AccountID  uint64             `gorm:"column:account_id;primaryKey;autoIncrement:true" json:"-"`
	RoleID     uint8              `gorm:"column:role_id;not null" json:"-"`
	Phone      *string            `gorm:"column:phone;size:20;unique" json:"phone" validate:"omitempty,vietnamese_phone"`
	Fullname   string             `gorm:"column:fullname;size:300" json:"fullname"`
	Password   *string            `gorm:"column:password;size:200" json:"password,omitempty" validate:"omitempty,min=8"`
	RePassword *string            `gorm:"-" json:"re_password" validate:"omitempty,eqfield=Password"`
	Images     []imagemodel.Image `gorm:"foreignKey:AccountID;references:AccountID" json:"images"`
	Email      *string            `gorm:"column:email;size:100;unique" json:"email" validate:"omitempty,email"`
	Version    int                `gorm:"default:1"`
	Gender     *common.Gender     `gorm:"column:gender;type:enum('Male', 'Female', 'Other')" json:"gender"`
	common.CommonFields
}

func (UpdateAccount) TableName() string {
	return Account{}.TableName()
}
