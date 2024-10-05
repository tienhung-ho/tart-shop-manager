package routerv1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	authhandler "tart-shop-manager/api/handler/auth"
	authmiddleware "tart-shop-manager/api/middleware/auth"
	accountv1 "tart-shop-manager/api/router/v1/account"
	categoryv1 "tart-shop-manager/api/router/v1/category"
	imagev1 "tart-shop-manager/api/router/v1/image"
	ingredientv1 "tart-shop-manager/api/router/v1/ingredient"
	orderv1 "tart-shop-manager/api/router/v1/order"
	productv1 "tart-shop-manager/api/router/v1/product"
	recipev1 "tart-shop-manager/api/router/v1/recipe"
	rolev1 "tart-shop-manager/api/router/v1/role"
	stockbatchv1 "tart-shop-manager/api/router/v1/stockbatch"
)

func NewRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	r.POST("/login", authhandler.LoginHandler(db))

	v1 := r.Group("/v1")
	v1.Use(authmiddleware.AuthRequire(db, rdb), authmiddleware.CasbinMiddleware())
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
		category := v1.Group("/category")
		{
			categoryv1.CategoryRouter(category, db, rdb)
		}
		order := v1.Group("/order")
		{
			orderv1.OrderRouter(order, db, rdb)
		}
		stockBatch := v1.Group("/stock-batch")
		{
			stockbatchv1.StockBatchRouter(stockBatch, db, rdb)
		}
		ingredient := v1.Group("/ingredient")
		{
			ingredientv1.IngredientRouter(ingredient, db, rdb)
		}
		image := v1.Group("/image")
		{
			imagev1.ImageRouter(image, db)
		}
		recipe := v1.Group("/recipe")
		{
			recipev1.RecipeRouter(recipe, db, rdb)
		}
	}
	return r
}
