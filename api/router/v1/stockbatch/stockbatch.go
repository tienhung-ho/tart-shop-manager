package stockbatchv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	stockbatchhandler "tart-shop-manager/api/handler/stockbatch"
)

func StockBatchRouter(stockBatch *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	stockBatch.POST("/", stockbatchhandler.CreateStockBatchHandler(db))
}
