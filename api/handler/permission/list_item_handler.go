package permissionhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	permissionstorage "tart-shop-manager/internal/repository/mysql/permission"
	permissioncachestorage "tart-shop-manager/internal/repository/redis/permission"
	permissionbusiness "tart-shop-manager/internal/service/permission"
)

func ListItemPermission(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		condition := map[string]interface{}{}

		paging := paggingcommon.Paging{
			Page:  1,
			Limit: 9999,
		}

		var filter commonfilter.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := permissionstorage.NewMySQLPermission(db)
		cache := permissioncachestorage.NewRdbStorage(rdb)
		biz := permissionbusiness.NewListItemPermissionBiz(store, cache)

		records, err := biz.ListItemPermission(c, condition, &paging, &filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewSuccesResponse(records, paging, filter))

	}
}
