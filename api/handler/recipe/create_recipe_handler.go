package recipehandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	recipestorage "tart-shop-manager/internal/repository/mysql/recipe"
	recipebusiness "tart-shop-manager/internal/service/recipe"
)

func CreateRecipeHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data recipemodel.CreateRecipe

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := recipestorage.NewMySQLRecipe(db)
		biz := recipebusiness.NewCreateRecipeBusiness(store)

		recipeID, err := biz.CreateRecipe(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recipeID, "create new recipe successfully"))

	}
}
