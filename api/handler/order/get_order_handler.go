package orderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	orderstorage "tart-shop-manager/internal/repository/mysql/order"
	ordercache "tart-shop-manager/internal/repository/redis/order"
	orderbusiness "tart-shop-manager/internal/service/order"
)

func GetOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderstorage.NewMySQLOrder(db)
		cache := ordercache.NewRdbStorage(rdb)
		biz := orderbusiness.NewGetOrderBiz(store, cache)

		record, err := biz.GetOrder(c, map[string]interface{}{"order_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get order successfully"))

	}
}
