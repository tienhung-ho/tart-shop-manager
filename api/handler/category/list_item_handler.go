package categoryhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorystorage "tart-shop-manager/internal/repository/mysql/category"
	categorycache "tart-shop-manager/internal/repository/redis/category"
	categorybusiness "tart-shop-manager/internal/service/category"
)

func ListCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		condition := map[string]interface{}{
			//"status": []string{"pending", "active", "inactive"},
		}

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

		store := categorystorage.NewMySQLCategory(db)
		cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewListItemCategoryBiz(store, cache)

		records, err := biz.ListItem(c.Request.Context(), condition, &paging, &filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewSuccesResponse(records, paging, filter))
	}
}
