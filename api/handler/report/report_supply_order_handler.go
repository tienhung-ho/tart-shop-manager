package reporthandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	reportstorage "tart-shop-manager/internal/repository/mysql/report"
	reportbusiness "tart-shop-manager/internal/service/report"
)

func ReportSupplyOrderHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		condition := map[string]interface{}{
			"status": []string{"pending", "active", "inactive"},
		}

		var paging paggingcommon.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		paging.Process()

		var filter commonfilter.Filter

		if err := c.ShouldBindQuery(&filter); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := reportstorage.NewMySQLOrder(db)
		biz := reportbusiness.NewReportSupplyOrderBusiness(store)
		record, err := biz.ReportSupplyOrder(c, condition, &paging, &filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "generate report successfully"))
	}
}
