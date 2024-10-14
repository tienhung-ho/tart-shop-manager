package supplyorderv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	supplyorderhandler "tart-shop-manager/api/handler/supplyOrder"
)

func SupplyOrderRouter(supplyOrder *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	supplyOrder.POST("/", supplyorderhandler.CreateSupplyOrderHandler(db, rdb))
}
