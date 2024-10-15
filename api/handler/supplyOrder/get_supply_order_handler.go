package supplyorderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	supplyorderstorage "tart-shop-manager/internal/repository/mysql/supplyOrder"
	supplyordercache "tart-shop-manager/internal/repository/redis/supplyOrder"
	supplyorderbusiness "tart-shop-manager/internal/service/supplyOrder"
)

func GetSupplyOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := supplyorderstorage.NewMySQLSupplyOrder(db)
		cache := supplyordercache.NewRdbStorage(rdb)
		biz := supplyorderbusiness.NewGetSupplyOrderBiz(store, cache)

		record, err := biz.GetSupplyOrder(c, map[string]interface{}{"supplyorder_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get supply order successfully"))
	}
}
