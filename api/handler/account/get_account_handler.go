package accounthandler

import (
	"errors"
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

func GetAccountHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		store := accountstorage.NewMySQLAccount(db)
		rdbAccount := accountrdbstorage.NewRdbStorage(rdb)
		biz := accountbusiness.NewGetAccountBiz(store, rdbAccount)

		record, err := biz.GetAccount(c.Request.Context(), map[string]interface{}{"account_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get account successfully"))

	}
}

func GetAccountWithAccessTokenHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		accountID, ok := c.Get("id")
		if !ok {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(errors.New("account id not found")))
			c.Abort()
			return
		}

		store := accountstorage.NewMySQLAccount(db)
		rdbAccount := accountrdbstorage.NewRdbStorage(rdb)
		biz := accountbusiness.NewGetAccountBiz(store, rdbAccount)
		cond := map[string]interface{}{
			"account_id": accountID,
			"status":     common.StatusActive,
		}
		record, err := biz.GetAccount(c, cond)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get account successfully"))

	}
}
