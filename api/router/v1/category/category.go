package categoryv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	categoryhandler "tart-shop-manager/api/handler/category"
)

func CategoryRouter(category *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	category.GET("/:id", categoryhandler.GetCategoryHandler(db, rdb))
}
