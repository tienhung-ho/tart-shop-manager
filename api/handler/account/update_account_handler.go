package accounthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
	accountstorage "tart-shop-manager/internal/repository/mysql/account"
	accountrdbstorage "tart-shop-manager/internal/repository/redis/account"
	accountbusiness "tart-shop-manager/internal/service/account"
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
			c.JSON(http.StatusBadRequest, common.ErrCanNotBindEntity("Account", err))
			return
		}

		validate := validator.New()

		validate.RegisterValidation("vietnamese_phone", func(fl validator.FieldLevel) bool {
			return validation.IsValidVietnamesePhoneNumber(fl.Field().String())
		})

		// Thực hiện validate
		if err := validate.Struct(data); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, common.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := accountstorage.NewMySQLAccount(db)
		cache := accountrdbstorage.NewRdbStorage(rdb)
		biz := accountbusiness.NewUpdateAccount(store, cache)

		updatedRecord, err := biz.UpdateAccount(c.Request.Context(), map[string]interface{}{"account_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(updatedRecord, "update account successfully"))

	}
}
