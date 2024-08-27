package accountv1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	accounthandler "tart-shop-manager/api/handler/account"
)

func AccountRouter(acc *gin.RouterGroup, db *gorm.DB) {
	acc.GET("/", func(c *gin.Context) {})
	acc.POST("/", accounthandler.CreateAccountHandler(db))
}
