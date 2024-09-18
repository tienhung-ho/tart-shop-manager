package ingredientv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	ingredienthandler "tart-shop-manager/api/handler/ingredient"
)

func IngredientRouter(ingre *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	ingre.GET("/:id", ingredienthandler.GetIngredientHandler(db, rdb))
	//ingre.POST("/", ingredienthandler.CreateIngredientHandler(db))
	ingre.PATCH("/:id", ingredienthandler.UpdateIngredientHandler(db, rdb))
}
