package stockbatchhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	stockbatchbusiness "tart-shop-manager/internal/service/stockbatch"
)

func CreateStockBatchHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data stockbatchmodel.CreateStockBatch

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		validate := validator.New()

		if err := validate.Struct(data); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := stockbatchstorage.NewMySQLStockBatch(db)
		biz := stockbatchbusiness.NewCreateStockBatchBiz(store)

		stockBatch, err := biz.CreateStockBatch(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(stockBatch, "create new stock batch successfully"))
	}
}
