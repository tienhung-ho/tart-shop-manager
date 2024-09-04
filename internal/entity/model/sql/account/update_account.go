package accountmodel

import "tart-shop-manager/internal/common"

type UpdateAccount struct {
	RoleID     uint8          `gorm:"column:role_id;not null" json:"role_id"`
	Phone      *string        `gorm:"column:phone;size:20;unique" json:"phone" validate:"omitempty,vietnamese_phone"`
	Fullname   string         `gorm:"column:fullname;size:300" json:"fullname"`
	AvatarURL  string         `gorm:"column:avatar_url;size:255" json:"avatar_url"`
	Password   *string        `gorm:"column:password;size:200" json:"password" validate:"omitempty,min=8"`
	RePassword *string        `gorm:"-" json:"re_password" validate:"required_if=Password omitempty,eqfield=Password"`
	Email      *string        `gorm:"column:email;size:100;unique" json:"email" validate:"omitempty,email"`
	Version    int            `gorm:"default:1"`
	Gender     *common.Gender `gorm:"column:gender;type:enum('Male', 'Female', 'Other')" json:"gender"`
	*common.CommonFields
}

func (UpdateAccount) TableName() string {
	return Account{}.TableName()
}
