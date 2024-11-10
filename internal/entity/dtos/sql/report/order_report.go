package reportmodel

import (
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

var (
	OrderReportName = "OrderReport"
)

type product struct {
	ProductID  uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	Name       string  `gorm:"column:name;size:200;not null" json:"name"`
	Price      float64 `gorm:"column:price;type:decimal(11,2)" json:"price"`
	ImageURL   string  `gorm:"column:image_url;size:300;not null" json:"image_url"`
	CategoryID uint64  `gorm:"column:category_id;not null" json:"category_id"`
	//Category        *categorymodel.Category `gorm:"foreignKey:CategoryID" json:"category"` // Liên kết với Category
	//*common.CommonFields
}

func (product) TableName() string {
	return "Product"
}

type Recipe struct {
	RecipeID  uint64   `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	ProductID uint64   `gorm:"column:product_id;not null" json:"product_id"`
	Product   *product `gorm:"foreignKey:ProductID;references:ProductID" json:"product"` // Liên kết với Product
	Size      string   `gorm:"column:size;type:enum('Small', 'Medium', 'Large');not null" json:"size"`
	Cost      float64  `gorm:"column:cost;not null" json:"cost"`
	//common.CommonFields
}

func (Recipe) TableName() string {
	return "Recipe"
}

type OrderItem struct {
	OrderID  uint64              `gorm:"column:order_id;primaryKey;autoIncrement:true" json:"order_id"`
	RecipeID uint64              `gorm:"column:recipe_id;primaryKey;autoIncrement" json:"recipe_id"`
	Quantity uint64              `gorm:"column:quantity;primaryKey;autoIncrement" json:"quantity"`
	Recipe   *recipemodel.Recipe `gorm:"foreignKey:RecipeID;references:RecipeID" json:"recipe,omitempty"`
	Price    float64             `gorm:"column:price;type:decimal(11,2)" json:"price"`
}

func (OrderItem) TableName() string {
	return "OrderRecipe"
}

type Order struct {
	OrderID     uint64  `gorm:"column:order_id;primaryKey;autoIncrement:true" json:"order_id"`
	AccountID   uint64  `gorm:"column:account_id;not null" json:"account_id"`
	TotalAmount float64 `gorm:"column:total_amount;type:decimal(11,2);not null;default:0.00" json:"total_amount"`
	Tax         float64 `gorm:"column:tax;type:decimal(10,2);default:0.00" json:"tax"`
	//Recipes     []recipemodel.Recipe `gorm:"many2many:OrderRecipe;foreignKey:OrderID;joinForeignKey:OrderID;References:RecipeID;joinReferences:RecipeID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"order_items,omitempty"`
	common.CommonFields
}

func (Order) TableName() string {
	return "Order"
}
