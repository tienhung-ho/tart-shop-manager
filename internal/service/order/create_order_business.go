package orderbusiness

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

type CreateOrderStorage interface {
	GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error)
	CreateOrder(ctx context.Context, data *ordermodel.CreateOrder) (uint64, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type CreateOrderItemStorage interface {
	CreateOrderItems(ctx context.Context, data []ordermodel.CreateOrderItem) error
}

type CreateOrderCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type createOrderBusiness struct {
	store           CreateOrderStorage
	cache           CreateOrderCache
	orderItemStore  CreateOrderItemStorage
	recipeStorage   RecipeStorage
	stockBatchStore StockBatchStorage
}

func NewCreateOrderBiz(
	store CreateOrderStorage,
	orderItemStore CreateOrderItemStorage,
	recipeStorage RecipeStorage,
	stockBatchStore StockBatchStorage,
	cache CreateOrderCache,
) *createOrderBusiness {
	return &createOrderBusiness{
		store,
		cache,
		orderItemStore,
		recipeStorage,
		stockBatchStore,
	}
}

func (biz *createOrderBusiness) CreateOrder(ctx context.Context, data *ordermodel.CreateOrder) (uint64, error) {

	var sizes []string
	var productIDs []uint64
	var profitProportion []float64
	var quantityProducts []uint64

	for _, item := range data.OrderItems {
		sizes = append(sizes, item.Size)
		productIDs = append(productIDs, item.ProductID)
		profitProportion = append(profitProportion, item.Price)
		quantityProducts = append(quantityProducts, item.Quantity)
	}

	cond := map[string]interface{}{}

	paging := paggingcommon.Paging{}
	paging.Process()
	filter := commonfilter.Filter{
		Recipe: commonfilter.Recipe{
			Sizes:      sizes,
			ProductIDs: productIDs,
		},
	}

	recipes, err := biz.recipeStorage.ListItem(ctx, cond, &paging, &filter)
	if err != nil {
		return 0, common.ErrCannotListEntity("Recipe", err)
	}

	var orderID uint64
	var totalAmount float64
	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		// Kiểm tra số lượng recipes và order items
		if len(recipes) != len(data.OrderItems) {
			return common.ErrInternal(errors.New("mismatch between recipes and order items"))
		}

		// Tính toán tổng số tiền
		for i, item := range recipes {
			totalAmount += float64(quantityProducts[i]) * (item.Cost + ((item.Cost * profitProportion[i]) / 100))
		}

		data.TotalAmount = totalAmount
		data.Tax = totalAmount * 10 / 100

		orderData := *data
		orderData.OrderItems = nil // Tách OrderItems ra

		// Tạo đơn hàng
		orderID, err = biz.store.CreateOrder(txCtx, &orderData)
		if err != nil {
			return common.ErrCannotCreateEntity(ordermodel.EntityName, err)
		}

		// Step 3: Extract Ingredients and Quantities
		ingredientMap := make(map[uint64]float64) // ingredientID -> totalQuantity
		for i, recipe := range recipes {
			for _, ingredient := range recipe.RecipeIngredients {
				ingredientMap[ingredient.IngredientID] += ingredient.Quantity * float64(quantityProducts[i])
			}
		}

		// Step 4: Fetch Stock Batches for All Ingredients
		var ingredientIDs []uint64
		for id := range ingredientMap {
			ingredientIDs = append(ingredientIDs, id)
		}

		stockBatches, err := biz.stockBatchStore.GetStockBatchesByIngredientIDs(txCtx, ingredientIDs)
		if err != nil {
			return common.ErrCannotGetEntity("StockBatch", err)
		}

		// Organize stock batches by ingredientID
		stockMap := make(map[uint64][]stockbatchmodel.StockBatch)
		for _, batch := range stockBatches {
			stockMap[batch.IngredientID] = append(stockMap[batch.IngredientID], batch)
		}

		var updateStockBatches []stockbatchmodel.UpdateStockBatch
		for ingredientID, totalQuantity := range ingredientMap {
			var totalStock float64

			for _, batch := range stockMap[ingredientID] {
				totalStock += batch.Quantity
			}

			if totalStock < totalQuantity {
				//log.Print(fmt.Sprintf("out of quantity at order %d", ingredientID))
				return common.ErrOutOffQuantity(stockbatchmodel.EntityName, errors.New(fmt.Sprintf("out of quantity at order %d", ingredientID)))
			}

			remainingQuantity := totalQuantity

			// Sắp xếp các stock batches theo ReceivedDate (FIFO)
			sort.Slice(stockMap[ingredientID], func(i, j int) bool {
				return stockMap[ingredientID][i].ReceivedDate.Before(stockMap[ingredientID][j].ReceivedDate)
			})

			for _, batch := range stockMap[ingredientID] {
				if remainingQuantity <= 0 {
					break
				}

				if batch.Quantity == 0 {
					// Bỏ qua stockBatch đã hết
					//log.Printf("Skipping StockBatchID: %d with Quantity: 0\n", batch.StockBatchID)
					continue
				}

				if batch.Quantity >= remainingQuantity {
					// Cập nhật trừ remainingQuantity từ batch này
					newQuantity := batch.Quantity - remainingQuantity
					updateStockBatches = append(updateStockBatches, stockbatchmodel.UpdateStockBatch{
						StockBatchID: batch.StockBatchID,
						IngredientID: ingredientID,
						Quantity:     &newQuantity,
					})
					remainingQuantity = 0
				} else {
					// Trừ toàn bộ số lượng của batch này và tiếp tục
					zero := 0.0
					updateStockBatches = append(updateStockBatches, stockbatchmodel.UpdateStockBatch{
						StockBatchID: batch.StockBatchID,
						IngredientID: ingredientID,
						Quantity:     &zero,
					})
					remainingQuantity -= batch.Quantity
				}
			}

			if remainingQuantity > 0 {
				return fmt.Errorf("số lượng stock không đủ cho ingredient ID %d sau khi trừ từ các stockbatch", ingredientID)
			}
		}

		// Cập nhật stock batches
		//_, err = biz.stockBatchStore.UpdateStockBatches(txCtx, cond, updateStockBatches)
		if err != nil {
			return common.ErrCannotUpdateEntity("StockBatch", err)
		}

		// Step 7: Create Order Items
		var orderItems []ordermodel.CreateOrderItem
		for i, item := range data.OrderItems {
			orderItem := ordermodel.CreateOrderItem{
				OrderID:  orderID,
				RecipeID: recipes[i].RecipeID,
				Quantity: item.Quantity,
				Price:    recipes[i].Cost + ((recipes[i].Cost * profitProportion[i]) / 100),
			}
			orderItems = append(orderItems, orderItem)
		}

		err = biz.orderItemStore.CreateOrderItems(txCtx, orderItems)
		if err != nil {
			return common.ErrCannotCreateEntity(ordermodel.EntityName, err)
		}

		if err = biz.cache.DeleteListCache(txCtx, ordermodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(ordermodel.EntityName, err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}
