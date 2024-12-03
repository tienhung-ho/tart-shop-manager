package permissionv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	permissionhandler "tart-shop-manager/api/handler/permission"
)

func PermissionRouter(permission *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	permission.GET("/list", permissionhandler.ListItemPermission(db, rdb))
}
