package accountv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	accounthandler "tart-shop-manager/api/handler/account"
	authhandler "tart-shop-manager/api/handler/auth"
	authmiddleware "tart-shop-manager/api/middleware/auth"
)

func AccountRouter(acc *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	acc.GET("/:id", authmiddleware.AuthRequire(), accounthandler.GetAccountHandler(db, rdb))
	acc.POST("/login", authhandler.LoginHandler(db))
	acc.POST("/", authmiddleware.AuthRequire(), accounthandler.CreateAccountHandler(db))
	acc.PATCH("/:id", authmiddleware.AuthRequire(), accounthandler.UpdateAccountHandler(db, rdb))
	acc.DELETE("/:id", authmiddleware.AuthRequire(), accounthandler.DeleteAccountHandler(db, rdb))
	acc.GET("/list", accounthandler.ListAccountHandler(db, rdb))
}
