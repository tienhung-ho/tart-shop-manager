package rolehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	rolecache "tart-shop-manager/internal/repository/redis/role"
	rolebusiness "tart-shop-manager/internal/service/role"
)

func DeleteRoleHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := rolestorage.NewMySQLRole(db)
		cache := rolecache.NewRdbStorage(rdb)
		biz := rolebusiness.NewDeleteRoleBiz(store, cache)

		if err := biz.DeleteRole(c.Request.Context(), map[string]interface{}{"role_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete role successfully"))
	}
}
