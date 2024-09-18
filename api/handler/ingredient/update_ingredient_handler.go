package ingredienthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	ingredientcache "tart-shop-manager/internal/repository/redis/ingredient"
	ingredientbusiness "tart-shop-manager/internal/service/ingredient"
)

func UpdateIngredientHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		var data ingredientmodel.UpdateIngredient

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := ingredientstorage.NewMySQLIngredient(db)
		cache := ingredientcache.NewRdbStorage(rdb)
		biz := ingredientbusiness.NewUpdateIngredientBiz(store, cache)

		newData, err := biz.UpdateIngredient(c, map[string]interface{}{"ingredient_id": id}, &data)

		if err != nil {
			
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(newData, "update ingredient successfully"))
	}
}
