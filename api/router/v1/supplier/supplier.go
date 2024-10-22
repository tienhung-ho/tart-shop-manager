package supplierv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	supplierhandler "tart-shop-manager/api/handler/supplier"
)

func SupplierRouter(supplier *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	supplier.GET("/:id", supplierhandler.GetSupplierHandler(db, rdb))
	supplier.POST("/", supplierhandler.CreateSupplierHandler(db, rdb))
	supplier.PATCH("/:id", supplierhandler.UpdateSupplierHandler(db, rdb))
	supplier.DELETE("/:id", supplierhandler.DeleteSupplierHandler(db, rdb))
	supplier.GET("/list", supplierhandler.ListItemHandler(db, rdb))
}
