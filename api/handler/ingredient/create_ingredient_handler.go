package ingredienthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	ingredientcache "tart-shop-manager/internal/repository/redis/ingredient"
	ingredientbusiness "tart-shop-manager/internal/service/ingredient"
)

func CreateIngredientHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data ingredientmodel.CreateIngredient
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

		store := ingredientstorage.NewMySQLIngredient(db)
		cache := ingredientcache.NewRdbStorage(rdb)
		biz := ingredientbusiness.NewCreateIngredientBiz(store, cache)

		recordID, err := biz.CreateIngredient(c, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordID, "create new ingredient successfully"))

	}
}
