package producthandler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	productstorage "tart-shop-manager/internal/repository/mysql/product"
	stockbatchstorage "tart-shop-manager/internal/repository/mysql/stockbatch"
	productcache "tart-shop-manager/internal/repository/redis/product"
	productbusiness "tart-shop-manager/internal/service/product"
)

func ListProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		condition := map[string]interface{}{
			//"status": []string{"pending", "active", "inactive"},
		}

		var paging paggingcommon.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		paging.Process()

		var filter commonfilter.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInternal(err))
			return
		}

		store := productstorage.NewMySQLProduct(db)
		cache := productcache.NewRdbStorage(rdb)
		stockbatchStore := stockbatchstorage.NewMySQLStockBatch(db)
		biz := productbusiness.NewListItemBiz(store, cache, stockbatchStore)

		records, err := biz.ListItem(c.Request.Context(), condition, &paging, &filter)

		if err != nil {
			fmt.Printf("Error sorting products: %v\n", err)
			// Return the error to the client
			var appErr *common.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, common.ErrCannotSort(productmodel.EntityName, err))
			}
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, common.NewSuccesResponse(records, paging, filter))

	}
}
