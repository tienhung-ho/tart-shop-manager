package stockbatchhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	stockbatchcache "tart-shop-manager/internal/repository/redis/stockBatch"
	stockbatchbusiness "tart-shop-manager/internal/service/stockbatch"
)

func UpdateStockBatchHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		var data stockbatchmodel.UpdateStockBatch

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrCanNotBindEntity(stockbatchmodel.EntityName, err))
			c.Abort()
			return
		}

		store := stockbatchstorage.NewMySQLStockBatch(db)
		cache := stockbatchcache.NewRdbStorage(rdb)
		biz := stockbatchbusiness.NewUpdateStockBatchBiz(store, cache)

		record, err := biz.UpdateStockBatch(c, map[string]interface{}{"stockbatch_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "updated stock batch successfully"))

	}
}
