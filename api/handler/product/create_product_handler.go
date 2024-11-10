package producthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	imagestorage "tart-shop-manager/internal/repository/mysql/image"
	productstorage "tart-shop-manager/internal/repository/mysql/product"
	productbusiness "tart-shop-manager/internal/service/product"
)

func CreateProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data productmodel.CreateProduct

		if err := c.ShouldBind(&data); err != nil {
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

		store := productstorage.NewMySQLProduct(db)
		cloud := imagestorage.NewMySQLImage(db)
		biz := productbusiness.NewCreateProductBusiness(store, cloud)

		recordId, err := biz.CreateProduct(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordId, "create product successfully"))
	}

}
