package reportv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	reporthandler "tart-shop-manager/api/handler/report"
)

func ReportRouter(report *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	report.GET("/order", reporthandler.GetReportHandler(db))
	report.GET("/supply-order", reporthandler.ReportSupplyOrderHandler(db))
}
