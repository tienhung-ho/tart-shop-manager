package accounthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	"tart-shop-manager/internal/entity/dtos/sql/account"
	"tart-shop-manager/internal/repository/mysql/account"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	accountrdbstorage "tart-shop-manager/internal/repository/redis/account"
	accountbusiness "tart-shop-manager/internal/service/account"
	casbinutil "tart-shop-manager/internal/util/policies"
	validation "tart-shop-manager/internal/validate"
)

func CreateAccountHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data accountmodel.CreateAccount

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		validate := validator.New()

		validate.RegisterValidation("vietnamese_phone", func(fl validator.FieldLevel) bool {
			return validation.IsValidVietnamesePhoneNumber(fl.Field().String())
		})

		// Thực hiện validate
		err := validate.Struct(data)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		enforcer := casbinutil.GetEnforcer()

		store := accountstorage.NewMySQLAccount(db)
		cache := accountrdbstorage.NewRdbStorage(rdb)
		roleStore := rolestorage.NewMySQLRole(db)
		auth := casbinutil.NewCasbinAuthorization(enforcer)
		biz := accountbusiness.NewCreateAccountbiz(store, cache, roleStore, auth)

		recordId, err := biz.CreateAccount(c, &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(recordId, "create account successfully"))
	}
}
