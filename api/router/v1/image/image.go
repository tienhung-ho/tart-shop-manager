package imagev1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	imagehandler "tart-shop-manager/api/handler/image"
	cloudmiddleware "tart-shop-manager/api/middleware/cloud"
)

func ImageRouter(image *gin.RouterGroup, db *gorm.DB) {
	image.POST("/", cloudmiddleware.ImageValidationMiddleware(), imagehandler.CreateImageHandler(db))
	image.DELETE("/:id", imagehandler.DeleteImageHandler(db))
}
