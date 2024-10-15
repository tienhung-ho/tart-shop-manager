package supplyorderbusiness

import (
	"context"
	"fmt"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	databaseutil "tart-shop-manager/internal/util/database"
)

type CreateSupplyOrderStorage interface {
	CreateSupplyOrder(ctx context.Context, data *supplyordermodel.CreateSupplyOrder) (uint64, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type createSupplyOrderBusiness struct {
	store           CreateSupplyOrderStorage
	storeItem       CreateSupplyOrderItemStorage
	storeIngredient IngredientStorage
	storeStockBatch StockBatchStorage
}

func NewCreateSupplyOrderBusiness(store CreateSupplyOrderStorage,
	storeItem CreateSupplyOrderItemStorage, storeIngredient IngredientStorage,
	storeStockBatch StockBatchStorage) *createSupplyOrderBusiness {
	return &createSupplyOrderBusiness{store, storeItem, storeIngredient,
		storeStockBatch}
}

func (biz *createSupplyOrderBusiness) CreateSupplyOrder(ctx context.Context, data *supplyordermodel.CreateSupplyOrder) (uint64, error) {
	var totalAmount float64 = 0.0
	ingredientIDs := make([]uint64, len(data.Ingredients))
	for i, ing := range data.Ingredients {
		ingredientIDs[i] = ing.IngredientID
		totalAmount += ing.Price
	}

	// Sử dụng hàm ListItem để lấy danh sách nguyên liệu tồn tại
	cond := map[string]interface{}{}
	paging := &paggingcommon.Paging{
		Page:  1,
		Limit: len(ingredientIDs),
	}
	filter := &commonfilter.Filter{
		IDs: ingredientIDs,
	}

	existingIngredients, err := biz.storeIngredient.ListItem(ctx, cond, paging, filter)
	if err != nil {
		return 0, common.ErrCannotGetEntity("ingredients", err)
	}

	// So sánh danh sách ingredientIDs và existingIngredients để tìm ra các ingredient_id không tồn tại
	existingIngredientIDs := make([]uint64, len(existingIngredients))
	for i, ing := range existingIngredients {
		existingIngredientIDs[i] = uint64(ing.IngredientID)
	}

	missingIngredientIDs := databaseutil.Difference(ingredientIDs, existingIngredientIDs)
	if len(missingIngredientIDs) > 0 {
		return 0, common.ErrInvalidRequest(fmt.Errorf("ingredients not found: %v", missingIngredientIDs))
	}

	var recordID uint64

	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		data.TotalAmount = totalAmount
		recordID, err = biz.store.CreateSupplyOrder(ctx, data)

		if err != nil {
			return common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
		}

		var stockBatches []stockbatchmodel.CreateStockBatch
		for _, item := range data.Ingredients {
			stockBatches = append(stockBatches, stockbatchmodel.CreateStockBatch{
				Quantity:       item.Quantity,
				ExpirationDate: item.ExpirationDate,
				ReceivedDate:   item.ReceivedDate,
				IngredientID:   uint(item.IngredientID),
			})
		}
		// Thực hiện bulk insert StockBatch
		stockIDs, err := biz.storeStockBatch.CreateStockBatches(txCtx, stockBatches)
		if err != nil {
			return common.ErrCannotCreateEntity(stockbatchmodel.EntityName, err)
		}

		// Chuẩn bị dữ liệu cho SupplyOrderItem
		var supOrderItems []supplyordermodel.CreateSupplyOrderItem
		for i, item := range data.Ingredients {
			supOrderItems = append(supOrderItems, supplyordermodel.CreateSupplyOrderItem{
				IngredientID:  item.IngredientID,
				Price:         item.Price,
				Quantity:      item.Quantity,
				Unit:          item.Unit,
				SupplyOrderID: recordID,
				StockBatchID:  stockIDs[i],
			})
		}

		err = biz.storeItem.CreateSupplyOrderItem(txCtx, supOrderItems)
		if err != nil {
			return common.ErrCannotCreateEntity("SupplyOrderItem", err)
		}

		return nil
	})

	if err != nil {
		return 0, common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
	}

	return recordID, nil
}
