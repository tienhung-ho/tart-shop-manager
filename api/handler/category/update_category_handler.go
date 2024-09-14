package categoryhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	categorystorage "tart-shop-manager/internal/repository/mysql/category"
	categorycache "tart-shop-manager/internal/repository/redis/category"
	categorybusiness "tart-shop-manager/internal/service/category"
)

func UpdateCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		var data categorymodel.UpdateCategory

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewUpdateCategoryBiz(store, cache)

		updatedCategory, err := biz.UpdateCategory(c, map[string]interface{}{"category_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(updatedCategory, "update category successfully"))

	}
}
