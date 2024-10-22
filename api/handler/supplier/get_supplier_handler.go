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

func GetSupplierHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := supplierstorage.NewMySQLSupplier(db)
		cache := suppliercache.NewRdbStorage(rdb)
		biz := supplierbusiness.NewGetSupplierBiz(store, cache)

		record, err := biz.GetSupplier(c, map[string]interface{}{"supplier_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get supplier successfully"))
	}
}
