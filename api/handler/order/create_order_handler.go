package orderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	orderstorage "tart-shop-manager/internal/repository/mysql/order"
	orderitemstorage "tart-shop-manager/internal/repository/mysql/orderItem"
	recipestorage "tart-shop-manager/internal/repository/mysql/recipe"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	ordercache "tart-shop-manager/internal/repository/redis/order"
	orderbusiness "tart-shop-manager/internal/service/order"
)

func CreateOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data ordermodel.CreateOrder

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderstorage.NewMySQLOrder(db)
		cache := ordercache.NewRdbStorage(rdb)
		orderItemStore := orderitemstorage.NewMySQLOrder(db)
		recipeStore := recipestorage.NewMySQLRecipe(db)
		stockbatchStore := stockbatchstorage.NewMySQLStockBatch(db)

		biz := orderbusiness.NewCreateOrderBiz(store, cache, orderItemStore, recipeStore, stockbatchStore)

		recordID, err := biz.CreateOrder(c, &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordID, "create new order successfully"))
	}
}
