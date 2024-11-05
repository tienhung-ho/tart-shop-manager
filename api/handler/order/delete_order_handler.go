package orderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	orderstorage "tart-shop-manager/internal/repository/mysql/order"
	orderitemstorage "tart-shop-manager/internal/repository/mysql/orderItem"
	orderbusiness "tart-shop-manager/internal/service/order"
)

func DeleteOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderstorage.NewMySQLOrder(db)
		orderItemStore := orderitemstorage.NewMySQLOrder(db)
		biz := orderbusiness.NewDeleteOrderBiz(store, orderItemStore)

		if err := biz.DeleteOrder(c, map[string]interface{}{"order_id": id}); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete order successfully"))

	}
}
