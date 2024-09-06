package rolehandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	permissionstorage "tart-shop-manager/internal/repository/mysql/permission"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	rolebusiness "tart-shop-manager/internal/service/role"
	casbinutil "tart-shop-manager/internal/util/policies"
)

func CreateRoleHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data rolemodel.CreateRole

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		// // Define model and policy paths
		modelPath := filepath.Join(cwd, "config/casbin", "rbac_model.conf")

		// Initialize Casbin Enforcers
		enforcer, err := casbinutil.InitEnforcer(db, modelPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			return
		}

		auth := casbinutil.NewCasbinAuthorization(enforcer)

		store := rolestorage.NewMySQLRole(db)
		perStore := permissionstorage.NewMySQLPermission(db)
		biz := rolebusiness.NewCreateRoleBiz(store, perStore, auth)

		recordId, err := biz.CreateRole(c.Request.Context(), &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordId, "create role successfully"))
	}
}
