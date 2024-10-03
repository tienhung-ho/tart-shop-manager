package imagemodel

type UpdateImage struct {
	ImageID    uint64  `json:"image_id" form:"-"`
	ProductID  *uint64 `gorm:"column:product_id;not null" json:"product_id"`
	CategoryID *uint64 `gorm:"column:category_id;index" json:"category_id,omitempty"` // Cho phép NULL
	AccountID  *uint64 `gorm:"column:account_id;index" json:"account_id,omitempty"`
}

func (UpdateImage) TableName() string {
	return Image{}.TableName()
}
