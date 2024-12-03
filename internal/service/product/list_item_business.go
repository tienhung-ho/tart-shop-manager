package productbusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error)
}

type ListItemCache interface {
	SaveProduct(ctx context.Context, data interface{}, morekeys ...string) error
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
	SaveFilter(ctx context.Context, filter *commonfilter.Filter, morekeys ...string) error
	ListItem(ctx context.Context, key string) ([]productmodel.Product, error)
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
	GetFilter(ctx context.Context, key string) (*commonfilter.Filter, error)
}

type StockBatchStorage interface {
	GetStockBatchesByIngredientIDs(ctx context.Context,
		ingredientIDs []uint64) ([]stockbatchmodel.StockBatch, error)
}

type listItemBusiness struct {
	store           ListItemStorage
	cache           ListItemCache
	stockBatchStore StockBatchStorage
}

func NewListItemBiz(store ListItemStorage, cache ListItemCache, stockBatchStore StockBatchStorage) *listItemBusiness {
	return &listItemBusiness{store, cache, stockBatchStore}
}

func (biz *listItemBusiness) ListItem(ctx context.Context, cond map[string]interface{},
	pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error) {
	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *pagging
	filterCopy := *filter

	// Tạo khóa cache
	baseKey, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", productmodel.EntityName),
	})

	productKey := baseKey
	pagingKey := baseKey + ":paging"

	if err != nil {
		return nil, common.ErrCannotGenerateKey(productmodel.EntityName, err)
	}

	// Gọi cache với khóa đã tạo
	records, err := biz.cache.ListItem(ctx, productKey)
	if err != nil {
		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	if len(records) != 0 {
		cachedPaging, err := biz.cache.GetPaging(ctx, pagingKey)
		if err == nil {
			pagging.Page = cachedPaging.Page
			pagging.Total = cachedPaging.Total
			pagging.Limit = cachedPaging.Limit
			pagging.Sort = cachedPaging.Sort
		}
		return records, nil
	}

	// Gọi store để lấy dữ liệu từ database
	records, err = biz.store.ListItem(ctx, cond, pagging, filter, morekeys...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
		}
		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	ingredientIDsMap := make(map[uint64]bool)
	for _, product := range records {
		for _, recipe := range product.Recipes {
			for _, ingredient := range recipe.RecipeIngredients {
				ingredientIDsMap[ingredient.IngredientID] = true
			}
		}
	}

	// Chuyển đổi map sang slice để lấy các ingredientID duy nhất
	var ingredientIDs []uint64
	for id := range ingredientIDsMap {
		ingredientIDs = append(ingredientIDs, id)
	}

	// Kiểm tra nếu không có ingredient IDs, không cần gọi GetStockBatches
	var stockMap map[uint64]float64
	if len(ingredientIDs) > 0 {
		// Step 2: Lấy StockBatch cho tất cả ingredientIDs
		stockBatches, err := biz.stockBatchStore.GetStockBatchesByIngredientIDs(ctx, ingredientIDs)
		if err != nil {
			return nil, common.ErrCannotGetEntity("StockBatch", err)
		}

		stockMap = make(map[uint64]float64)
		for _, stock := range stockBatches {
			stockMap[stock.IngredientID] += stock.Quantity
		}
	} else {
		// Nếu không có ingredient IDs, khởi tạo stockMap rỗng
		stockMap = make(map[uint64]float64)
	}

	for i := range records {
		product := &records[i]

		if len(product.Recipes) == 0 {
			// Nếu không có recipe nào, đặt available_in_stock là false
			product.AvailableInStock = false
			continue
		}

		productAvailable := true
		for _, recipe := range product.Recipes {
			for _, ingredient := range recipe.RecipeIngredients {
				neededQuantity := ingredient.Quantity
				availableQuantity, exists := stockMap[ingredient.IngredientID]

				if !exists || availableQuantity < neededQuantity {
					productAvailable = false
					break
				}
			}
			if !productAvailable {
				break
			}
		}

		// Cập nhật trường available_in_stock
		product.AvailableInStock = productAvailable
	}

	// Lưu vào cache với cùng khóa
	if len(records) != 0 {
		if err := biz.cache.SaveProduct(ctx, records, productKey); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
		if err := biz.cache.SavePaging(ctx, pagging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
	}

	return records, nil
}
