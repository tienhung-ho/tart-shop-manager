package accounthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	accountstorage "tart-shop-manager/internal/repository/mysql/account"
	imagestorage "tart-shop-manager/internal/repository/mysql/image"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	accountrdbstorage "tart-shop-manager/internal/repository/redis/account"
	accountbusiness "tart-shop-manager/internal/service/account"
	casbinutil "tart-shop-manager/internal/util/policies"
	validation "tart-shop-manager/internal/validate"
)

func UpdateAccountHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		var data accountmodel.UpdateAccount

		if err := c.ShouldBindJSON(&data); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, common.ErrCanNotBindEntity("Account", err))
			return
		}
		log.Print(data.Password, data.RePassword)

		validate := validator.New()

		validate.RegisterValidation("vietnamese_phone", func(fl validator.FieldLevel) bool {
			return validation.IsValidVietnamesePhoneNumber(fl.Field().String())
		})

		// Thực hiện validate
		if err := validate.Struct(data); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				log.Print(validationErrors)
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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

		store := accountstorage.NewMySQLAccount(db)
		cache := accountrdbstorage.NewRdbStorage(rdb)
		roleStore := rolestorage.NewMySQLRole(db)
		image := imagestorage.NewMySQLImage(db)
		auth := casbinutil.NewCasbinAuthorization(enforcer)
		biz := accountbusiness.NewUpdateAccount(store, roleStore, cache, image, auth)

		updatedRecord, err := biz.UpdateAccount(c, map[string]interface{}{"account_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		updatedRecord.Password = ""
		c.JSON(http.StatusOK, common.NewDataResponse(updatedRecord, "update account successfully"))

	}
}
