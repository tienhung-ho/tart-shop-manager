package supplierhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	supplierstorage "tart-shop-manager/internal/repository/mysql/supplier"
	suppliercache "tart-shop-manager/internal/repository/redis/supplier"
	supplierbusiness "tart-shop-manager/internal/service/supplier"
)

func DeleteSupplierHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := supplierstorage.NewMySQLSupplier(db)
		cache := suppliercache.NewRdbStorage(rdb)
		biz := supplierbusiness.NewDeleteSupplierBiz(store, cache)

		if err := biz.DeleteSupplier(c, map[string]interface{}{"supplier_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete supplier successfully"))
	}
}
