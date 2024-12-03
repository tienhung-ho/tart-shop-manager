package recipehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	recipestorage "tart-shop-manager/internal/repository/mysql/recipe"
	recipeingredientstorage "tart-shop-manager/internal/repository/mysql/recipeIngredient"
	recipecache "tart-shop-manager/internal/repository/redis/recipe"
	recipebusiness "tart-shop-manager/internal/service/recipe"
)

func CreateRecipeHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data recipemodel.CreateRecipe

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := recipestorage.NewMySQLRecipe(db)
		ingredientStore := ingredientstorage.NewMySQLIngredient(db)
		cache := recipecache.NewRdbStorage(rdb)
		recipeIngredientStore := recipeingredientstorage.NewMySQLRecipeIngredient(db)

		biz := recipebusiness.NewCreateRecipeBusiness(store, cache, ingredientStore, recipeIngredientStore)

		recipeID, err := biz.CreateRecipe(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recipeID, "create new recipe successfully"))

	}
}
