package accountv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	accounthandler "tart-shop-manager/api/handler/account"
)

func AccountRouter(acc *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {

	acc.GET("/:id", accounthandler.GetAccountHandler(db, rdb))
	//acc.POST("/login", authhandler.LoginHandler(db))
	acc.POST("/", accounthandler.CreateAccountHandler(db))
	acc.PATCH("/:id", accounthandler.UpdateAccountHandler(db, rdb))
	acc.DELETE("/:id", accounthandler.DeleteAccountHandler(db, rdb))
	acc.GET("/list", accounthandler.ListAccountHandler(db, rdb))
}
