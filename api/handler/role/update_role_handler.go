package rolehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	permissionstorage "tart-shop-manager/internal/repository/mysql/permission"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	rolecache "tart-shop-manager/internal/repository/redis/role"
	rolebusiness "tart-shop-manager/internal/service/role"
	casbinutil "tart-shop-manager/internal/util/policies"
)

func UpdateRoleHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		var updateData rolemodel.UpdateRole

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		enforcer := casbinutil.GetEnforcer()

		store := rolestorage.NewMySQLRole(db)
		cache := rolecache.NewRdbStorage(rdb)
		perStore := permissionstorage.NewMySQLPermission(db)
		auth := casbinutil.NewCasbinAuthorization(enforcer)
		biz := rolebusiness.NewUpdateRolebiz(store, cache, perStore, auth)

		if err := biz.UpdateRole(c.Request.Context(), map[string]interface{}{"role_id": id}, &updateData); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "update role successfully"))

	}
}
