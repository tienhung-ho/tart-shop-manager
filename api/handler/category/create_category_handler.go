package categoryhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	categorystorage "tart-shop-manager/internal/repository/mysql/category"
	imagestorage "tart-shop-manager/internal/repository/mysql/image"
	categorybusiness "tart-shop-manager/internal/service/category"
)

func CreateCategoryHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data categorymodel.CreateCategory
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}
		validate := validator.New()

		err := validate.Struct(&data)
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

		store := categorystorage.NewMySQLCategory(db)
		cloud := imagestorage.NewMySQLImage(db)
		biz := categorybusiness.NewCreateCategoryBusiness(store, cloud)

		recordId, err := biz.CreateCategory(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordId, "create category successfully"))
	}
}
