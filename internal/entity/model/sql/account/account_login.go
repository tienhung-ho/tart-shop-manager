package accountmodel

type LoginAccount struct {
	Email    string `gorm:"column:email;size:100;not null;unique" json:"email" form:"email"`
	Password string `gorm:"column:password;size:100;not null" json:"password" form:"password"`
}

func (LoginAccount) TableName() string {
	return Account{}.TableName()
}
