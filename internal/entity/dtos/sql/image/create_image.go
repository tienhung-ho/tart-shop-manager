package imagemodel

type CreateImage struct {
	ImageID uint64 `gorm:"column:image_id;primaryKey;autoIncrement" json:"-" form:"-"`
	URL     string `gorm:"column:url;size:300;not null" json:"url" form:"url"`
	AltText string `gorm:"column:alt_text;size:255" json:"alt_text,omitempty" form:"alt_text"`
}

func (CreateImage) TableName() string {
	return Image{}.TableName()
}
