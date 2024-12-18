package recipev1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	recipehandler "tart-shop-manager/api/handler/recipe"
)

func RecipeRouter(recipe *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	recipe.GET("/:id", recipehandler.GetRecipeHandler(db, rdb))
	recipe.POST("/", recipehandler.CreateRecipeHandler(db, rdb))
	recipe.PATCH("/:id", recipehandler.UpdateRecipeHandler(db, rdb))
	recipe.DELETE("/:id", recipehandler.DeleteRecipeHandler(db, rdb))
	recipe.GET("/list", recipehandler.ListItemRecipe(db, rdb))
}
