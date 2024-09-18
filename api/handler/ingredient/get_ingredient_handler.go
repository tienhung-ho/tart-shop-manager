package ingredienthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	ingredientcache "tart-shop-manager/internal/repository/redis/ingredient"
	ingredientbusiness "tart-shop-manager/internal/service/ingredient"
)

func GetIngredientHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := ingredientstorage.NewMySQLIngredient(db)
		cache := ingredientcache.NewRdbStorage(rdb)
		biz := ingredientbusiness.NewGetIngredientBiz(store, cache)

		record, err := biz.GetIngredient(c, map[string]interface{}{"ingredient_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get ingredient successfully"))
	}
}
