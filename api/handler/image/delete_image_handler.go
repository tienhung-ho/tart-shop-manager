package imagehandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tart-shop-manager/internal/common"
	imagestorage "tart-shop-manager/internal/repository/mysql/image"
	imagebusiness "tart-shop-manager/internal/service/image"
)

func DeleteImageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			c.Abort()
			return
		}

		store := imagestorage.NewMySQLImage(db)
		biz := imagebusiness.NewDeleteImageBiz(store)

		if err := biz.DeleteImage(c, map[string]interface{}{"image_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewDataResponse(true, "delete image successfully"))
	}
}
