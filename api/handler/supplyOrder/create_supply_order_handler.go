package supplyorderhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	supplyorderstorage "tart-shop-manager/internal/repository/mysql/supplyOrder"
	supplyorderitemstorage "tart-shop-manager/internal/repository/mysql/supplyOrderItem"
	supplyorderbusiness "tart-shop-manager/internal/service/supplyOrder"
)

func CreateSupplyOrderHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data supplyordermodel.CreateSupplyOrder

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		validate := validator.New()
		
		// Thực hiện validate
		err := validate.Struct(data)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := supplyorderstorage.NewMySQLSupplyOrder(db)
		storeItem := supplyorderitemstorage.NewMySQLSupplyOrderItem(db)
		storeIngredient := ingredientstorage.NewMySQLIngredient(db)
		storeStockBathc := stockbatchstorage.NewMySQLStockBatch(db)
		biz := supplyorderbusiness.NewCreateSupplyOrderBusiness(store, storeItem, storeIngredient, storeStockBathc)

		recordID, err := biz.CreateSupplyOrder(c, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordID, "create supply order successfully"))
	}
}
