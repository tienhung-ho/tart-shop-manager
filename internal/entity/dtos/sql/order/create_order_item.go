package ordermodel

type CreateOrderItem struct {
	OrderID  uint64  `gorm:"column:order_id;not null" json:"order_id"`
	RecipeID uint64  `gorm:"column:recipe_id;not null" json:"recipe_id"`
	Quantity uint64  `gorm:"column:quantity;primaryKey;autoIncrement" json:"quantity"`
	Price    float64 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	//
	ProductID   uint64 `gorm:"-" json:"product_id"`
	Name        string `gorm:"-" json:"name"`
	Description string `gorm:"-" json:"description"`
	Size        string `gorm:"-" json:"size"`
}

func (*CreateOrderItem) TableName() string {
	return OrderItem{}.TableName()
}
