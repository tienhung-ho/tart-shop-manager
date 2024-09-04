package accounthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	accountstorage "tart-shop-manager/internal/repository/mysql/account"
	accountrdbstorage "tart-shop-manager/internal/repository/redis/account"
	accountbusiness "tart-shop-manager/internal/service/account"
)

func DeleteAccountHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := accountstorage.NewMySQLAccount(db)
		cache := accountrdbstorage.NewRdbStorage(rdb)
		biz := accountbusiness.NewDeleteAccountBiz(store, cache)

		if err := biz.DeleteAccount(c.Request.Context(), map[string]interface{}{"account_id": idInt}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "deleted successfully"))

	}
}
