package supplyorderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	supplyorderstorage "tart-shop-manager/internal/repository/mysql/supplyOrder"
	supplyorderitemstorage "tart-shop-manager/internal/repository/mysql/supplyOrderItem"
	supplyordercache "tart-shop-manager/internal/repository/redis/supplyOrder"
	supplyorderbusiness "tart-shop-manager/internal/service/supplyOrder"
)

func DeleteSupplyOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := supplyorderstorage.NewMySQLSupplyOrder(db)
		cache := supplyordercache.NewRdbStorage(rdb)
		storeItem := supplyorderitemstorage.NewMySQLSupplyOrderItem(db)
		storeIngredient := ingredientstorage.NewMySQLIngredient(db)
		storeStockBathc := stockbatchstorage.NewMySQLStockBatch(db)
		biz := supplyorderbusiness.NewDeleteSupplyOrderBiz(store, cache, storeItem, storeIngredient, storeStockBathc)

		if err := biz.DeleteSupplyOrder(c, map[string]interface{}{"supplyorder_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete supply order successfully"))

	}
}
