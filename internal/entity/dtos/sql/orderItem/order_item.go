package orderitemmodel

import recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"

var (
	EntityNameOrderItem = "OrderItem"
)

type OrderItem struct {
	OrderID  uint64              `gorm:"column:order_id;primaryKey;autoIncrement:true" json:"order_id"`
	RecipeID uint64              `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	Quantity uint64              `gorm:"column:quantity;primaryKey;autoIncrement" json:"quantity"`
	Recipe   *recipemodel.Recipe `gorm:"foreignKey:RecipeID;references:RecipeID" json:"recipe,omitempty"`
	Price    float64             `gorm:"column:price;type:decimal(11,2)" json:"price"`
	//
	ProductID   uint64 `gorm:"-" json:"product_id,omitempty"`
	Name        string `gorm:"-" json:"name,omitempty"`
	Description string `gorm:"-" json:"description,omitempty"`
	Size        string `gorm:"-" json:"size,omitempty"`
}

func (OrderItem) TableName() string {
	return "OrderRecipe"
}
