package recipehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	recipestorage "tart-shop-manager/internal/repository/mysql/recipe"
	recipecache "tart-shop-manager/internal/repository/redis/recipe"
	recipebusiness "tart-shop-manager/internal/service/recipe"
)

func GetRecipeHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := recipestorage.NewMySQLRecipe(db)
		cache := recipecache.NewRdbStorage(rdb)
		biz := recipebusiness.NewGetRecipeBiz(store, cache)

		record, err := biz.GetRecipe(c, map[string]interface{}{"recipe_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get recipe successfully"))
	}
}
