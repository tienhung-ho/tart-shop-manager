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

func GetStockBatchHandler(db *gorm.DB, rdb *redis.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := stockbatchstorage.NewMySQLStockBatch(db)
		cache := stockbatchcache.NewRdbStorage(rdb)
		biz := stockbatchbusiness.NewGetStockBatchBiz(store, cache)

		record, err := biz.GetStockBatch(c, map[string]interface{}{"stockbatch_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get stock batch successfully"))
	}
}
