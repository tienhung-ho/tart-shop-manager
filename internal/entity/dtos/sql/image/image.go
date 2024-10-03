package imagemodel

var (
	EntityName = "Image"
)

type Image struct {
	ImageID    uint64  `gorm:"column:image_id;primaryKey;autoIncrement" json:"image_id" form:"-"`
	URL        string  `gorm:"column:url;size:300;not null" json:"url" form:"image_url"`
	AltText    string  `gorm:"column:alt_text;size:255" json:"alt_text,omitempty" form:"alt_text"`
	ProductID  *uint64 `gorm:"column:product_id;index" json:"product_id,omitempty"`   // Cho phép NULL
	CategoryID *uint64 `gorm:"column:category_id;index" json:"category_id,omitempty"` // Cho phép NULL
	AccountID  *uint64 `gorm:"column:account_id;index" json:"account_id,omitempty"`
}

func (Image) TableName() string {
	return "Image"
}
