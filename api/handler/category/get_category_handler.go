package categoryhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	categorystorage "tart-shop-manager/internal/repository/mysql/category"
	categorycache "tart-shop-manager/internal/repository/redis/category"
	categorybusiness "tart-shop-manager/internal/service/category"
)

func GetCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewGetCategoryBiz(store, cache)

		record, err := biz.GetCategory(c.Request.Context(), map[string]interface{}{"category_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(record, "get account successfully"))

	}
}
