package ingredienthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientstorage "tart-shop-manager/internal/repository/mysql/ingredient"
	ingredientcache "tart-shop-manager/internal/repository/redis/ingredient"
	ingredientbusiness "tart-shop-manager/internal/service/ingredient"
)

func ListItemIngredient(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		condition := map[string]interface{}{}

		var paging paggingcommon.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		paging.Process()

		var filter commonfilter.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := ingredientstorage.NewMySQLIngredient(db)
		cache := ingredientcache.NewRdbStorage(rdb)
		biz := ingredientbusiness.NewListItemIngredientBiz(store, cache)

		records, err := biz.ListItem(c, condition, &paging, &filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewSuccesResponse(records, paging, filter))

	}
}
