package supplierhandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
	supplierstorage "tart-shop-manager/internal/repository/mysql/supplier"
	suppliercache "tart-shop-manager/internal/repository/redis/supplier"
	supplierbusiness "tart-shop-manager/internal/service/supplier"
)

func CreateSupplierHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data suppliermodel.CreateSupplier

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		varlidate := validator.New()

		if err := varlidate.Struct(&data); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := supplierstorage.NewMySQLSupplier(db)
		cache := suppliercache.NewRdbStorage(rdb)
		biz := supplierbusiness.NewCreateSupplierBusiness(store, cache)

		recordID, err := biz.CreateSupplier(c, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordID, "create supplier successfully"))
	}
}
