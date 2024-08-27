package accountmodel

import "tart-shop-manager/internal/common"

type CreateAccount struct {
	AccountID  uint64         `gorm:"column:account_id;primaryKey;autoIncrement:true" json:"-"`
	RoleID     uint8          `gorm:"column:role_id;not null" json:"role_id"`
	Phone      string         `gorm:"column:phone;size:20;not null;unique" json:"phone" validate:"required,vietnamese_phone"`
	Fullname   string         `gorm:"column:fullname;size:300" json:"fullname"`
	AvatarURL  string         `gorm:"column:avatar_url;size:255" json:"avatar_url"`
	Password   string         `gorm:"column:password;size:200;not null" json:"password" validate:"required,min=8"`
	RePassword string         `gorm:"-" json:"re_password" validate:"required,eqfield=Password"`
	Email      string         `gorm:"column:email;size:100;not null;unique" json:"email" validate:"required,email"`
	Gender     *common.Gender `gorm:"column:gender;type:enum('Male', 'Female', 'Other')" json:"gender"`
	*common.CommonFields
}

func (CreateAccount) TableName() string {
	return Account{}.TableName()
}
