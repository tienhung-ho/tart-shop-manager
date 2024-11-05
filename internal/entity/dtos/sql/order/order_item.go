package ordermodel

import (
	"gorm.io/gorm"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

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
	ProductID   uint64         `gorm:"-" json:"product_id"`
	Name        string         `gorm:"-" json:"name"`
	Description string         `gorm:"-" json:"description"`
	Size        string         `gorm:"-" json:"size"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OrderItem) TableName() string {
	return "OrderRecipe"
}
