package account

import (
	"tart-shop-manager/internal/common"
)

type CreateAccountRdb struct {
	AccountID  uint64         `json:"account_id"`
	RoleID     uint8          `json:"role_id"`
	Phone      string         `json:"phone" validate:"required,vietnamese_phone"`
	Fullname   string         `json:"fullname"`
	AvatarURL  string         `json:"avatar_url"`
	Password   string         `json:"password" validate:"required,min=8"`
	RePassword string         `json:"re_password" validate:"required,eqfield=Password"`
	Email      string         `json:"email" validate:"required,email"`
	Gender     *common.Gender `json:"gender"`
	*common.CommonFields
}
