package producthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	productstorage "tart-shop-manager/internal/repository/mysql/product"
	productcache "tart-shop-manager/internal/repository/redis/product"
	productbusiness "tart-shop-manager/internal/service/product"
)

func UpdateProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		var data productmodel.UpdateProduct

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
		}

		store := productstorage.NewMySQLProduct(db)
		cache := productcache.NewRdbStorage(rdb)
		biz := productbusiness.NewUpdatePruductBiz(store, cache)

		updatedProduct, err := biz.UpdateProduct(c.Request.Context(), map[string]interface{}{"product_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(updatedProduct, "update product successfully"))
	}
}