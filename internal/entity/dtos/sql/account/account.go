package accountmodel

import (
	"tart-shop-manager/internal/common"
	accountrdbmodel "tart-shop-manager/internal/entity/dtos/redis/account"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

const (
	EntityName = "Account"
)

var (
	SelectFields = []string{
		"account_id", // AccountID
		"role_id",    // RoleID
		"phone",      // Phone
		"fullname",   // Fullname
		"status",     // Status
		"email",      // Email
		"gender",     // Gender
		"password",
	}
)

type Account struct {
	AccountID uint64             `gorm:"column:account_id;primaryKey;autoIncrement" json:"account_id"`
	RoleID    uint8              `gorm:"column:role_id;not null" json:"role_id"`
	Role      *rolemodel.Role    `gorm:"foreignKey:RoleID;references:RoleID" json:"role"`
	Phone     string             `gorm:"column:phone;size:20;not null;unique" json:"phone"`
	Fullname  string             `gorm:"column:fullname;size:300" json:"fullname"`
	Password  string             `gorm:"column:password;size:200;not null" json:"password"`
	Email     string             `gorm:"column:email;size:100;not null;unique" json:"email"`
	Images    []imagemodel.Image `gorm:"foreignKey:AccountID;references:AccountID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE;" json:"images"`
	Gender    *common.Gender     `gorm:"column:gender;type:enum('Male', 'Female', 'Other')" json:"gender"`
	*common.CommonFields
}

func (Account) TableName() string {
	return "Account"
}

func (a Account) ToCreateAccount() *accountrdbmodel.CreateAccountRdb {
	return &accountrdbmodel.CreateAccountRdb{
		AccountID: a.AccountID,
		RoleID:    a.RoleID,
		Phone:     a.Phone,
		Fullname:  a.Fullname,
		Password:  a.Password,
		Email:     a.Email,
		Gender:    a.Gender,
		CommonFields: &common.CommonFields{
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			Status:    a.Status,
			CreatedBy: a.CreatedBy,
			UpdatedBy: a.UpdatedBy,
		},
	}
}

func (a Account) ToSimpleAccount() *Account {
	return &Account{
		AccountID: a.AccountID,
		RoleID:    a.RoleID,
		Phone:     a.Phone,
		Fullname:  a.Fullname,
		Role:      a.Role,
		Email:     a.Email,
		Gender:    a.Gender,
		CommonFields: &common.CommonFields{
			Status: a.Status,
		},
	}
}
