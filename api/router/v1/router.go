package routerv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	accountv1 "tart-shop-manager/api/router/v1/account"
	productv1 "tart-shop-manager/api/router/v1/product"
	rolev1 "tart-shop-manager/api/router/v1/role"
)

func NewRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		account := v1.Group("/account")
		{
			accountv1.AccountRouter(account, db, rdb)
		}

		role := v1.Group("/role")
		{
			rolev1.RoleRouter(role, db, rdb)
		}

		product := v1.Group("/product")
		{
			productv1.ProductRouter(product, db, rdb)
		}
	}
	return r
}
