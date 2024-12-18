package recipehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	recipestorage "tart-shop-manager/internal/repository/mysql/recipe"
	recipeingredientstorage "tart-shop-manager/internal/repository/mysql/recipeIngredient"
	recipecache "tart-shop-manager/internal/repository/redis/recipe"
	recipebusiness "tart-shop-manager/internal/service/recipe"
)

func DeleteRecipeHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := recipestorage.NewMySQLRecipe(db)
		cache := recipecache.NewRdbStorage(rdb)
		recipeIngredientStore := recipeingredientstorage.NewMySQLRecipeIngredient(db)
		biz := recipebusiness.NewDeleteRecipeBiz(store, cache, recipeIngredientStore)

		if err := biz.DeleteRecipe(c, map[string]interface{}{"recipe_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete recipe successfully"))
	}
}
