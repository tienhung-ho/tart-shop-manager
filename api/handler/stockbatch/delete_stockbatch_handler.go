package stockbatchhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	stockbatchcache "tart-shop-manager/internal/repository/redis/stockBatch"
	stockbatchbusiness "tart-shop-manager/internal/service/stockbatch"
)

func DeleteStockBatchHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := stockbatchstorage.NewMySQLStockBatch(db)
		cache := stockbatchcache.NewRdbStorage(rdb)
		biz := stockbatchbusiness.NewDeleteStockBatchBiz(store, cache)

		if err := biz.DeleteStockBatch(c, map[string]interface{}{"stockbatch_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete stock batch successfully"))
	}
}
