package orderv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	orderhandler "tart-shop-manager/api/handler/order"
)

func OrderRouter(order *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	order.GET("/:id", orderhandler.GetOrderHandler(db, rdb))
}
