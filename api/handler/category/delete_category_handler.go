package categoryhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	categorystorage "tart-shop-manager/internal/repository/mysql/category"
	imagestorage "tart-shop-manager/internal/repository/mysql/image"
	categorycache "tart-shop-manager/internal/repository/redis/category"
	categorybusiness "tart-shop-manager/internal/service/category"
)

func DeleteCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		cache := categorycache.NewRdbStorage(rdb)
		image := imagestorage.NewMySQLImage(db)
		biz := categorybusiness.NewDeleteCategoryBiz(store, cache, image)

		if err := biz.DeleteCategory(c, map[string]interface{}{"category_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "deleted category successfully"))

	}
}
