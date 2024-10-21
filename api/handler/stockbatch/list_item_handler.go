package stockbatchhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	stockbatchcache "tart-shop-manager/internal/repository/redis/stockBatch"
	stockbatchbusiness "tart-shop-manager/internal/service/stockbatch"
)

func ListItemStockBatchHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		condition := map[string]interface{}{
			"status": []string{"pending", "active", "inactive"},
		}

		var paging paggingcommon.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		paging.Process()

		var filter commonfilter.Filter

		if err := c.ShouldBindQuery(&filter); err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := stockbatchstorage.NewMySQLStockBatch(db)
		cache := stockbatchcache.NewRdbStorage(rdb)
		biz := stockbatchbusiness.NewListItemStockBatchBiz(store, cache)
		records, err := biz.ListItem(c, condition, &paging, &filter)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewSuccesResponse(records, paging, filter))
	}
}
