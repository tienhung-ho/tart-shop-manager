package rolev1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	rolehandler "tart-shop-manager/api/handler/role"
)

func RoleRouter(role *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	role.GET("/:id", rolehandler.GetRoleHandler(db, rdb))
	role.POST("/", rolehandler.CreateRoleHandler(db))
}
