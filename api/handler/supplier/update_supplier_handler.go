package supplierhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	supplierstorage "tart-shop-manager/internal/repository/mysql/supplier"
	suppliercache "tart-shop-manager/internal/repository/redis/supplier"
	supplierbusiness "tart-shop-manager/internal/service/supplier"
)

func UpdateSupplierHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		var data suppliermodel.UpdateSupplier
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrCanNotBindEntity(suppliermodel.EntityName, err))
			c.Abort()
			return
		}

		store := supplierstorage.NewMySQLSupplier(db)
		cache := suppliercache.NewRdbStorage(rdb)
		biz := supplierbusiness.NewUpdateSupplierBiz(store, cache)

		record, err := biz.UpdateSupplier(c, map[string]interface{}{"supplier_id": id}, &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "update supplier successfully"))
	}
}
