package account

import (
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

type CreateAccountRdb struct {
	AccountID  uint64             `json:"account_id"`
	RoleID     uint8              `json:"role_id"`
	Role       *rolemodel.Role    `json:"role"`
	Phone      string             `json:"phone" validate:"required,vietnamese_phone"`
	Fullname   string             `json:"fullname"`
	Password   string             `json:"password" validate:"required,min=8"`
	RePassword string             `json:"re_password" validate:"required,eqfield=Password"`
	Images     []imagemodel.Image `gorm:"foreignKey:AccountID;references:AccountID" json:"images"`
	Email      string             `json:"email" validate:"required,email"`
	Gender     *common.Gender     `json:"gender"`
	*common.CommonFields
}
