package supplyorderbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateSupplyOrderStorage interface {
	UpdateSupplyOrder(ctx context.Context, cond map[string]interface{}, data *supplyordermodel.UpdateSupplyOrder) (*supplyordermodel.SupplyOrder, error)
	GetSupplyOrder(ctx context.Context, cond map[string]interface{}) (*supplyordermodel.SupplyOrder, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UpdateSupplyOrderCache interface {
	SaveSupplyOrder(ctx context.Context, data interface{}, morekeys ...string) error
	DeleteSupplyOrder(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type updateSupplyOrderBusiness struct {
	store           UpdateSupplyOrderStorage
	cache           UpdateSupplyOrderCache
	storeItem       CreateSupplyOrderItemStorage
	storeIngredient IngredientStorage
	storeStockBatch StockBatchStorage
}

func NewUpdateSupplyOrderBiz(store UpdateSupplyOrderStorage, cache UpdateSupplyOrderCache,
	storeItem CreateSupplyOrderItemStorage, storeIngredient IngredientStorage,
	storeStockBatch StockBatchStorage) *updateSupplyOrderBusiness {
	return &updateSupplyOrderBusiness{store: store, cache: cache, storeItem: storeItem,
		storeStockBatch: storeStockBatch}
}

func (biz *updateSupplyOrderBusiness) UpdateSupplyOrder(ctx context.Context, cond map[string]interface{},
	data *supplyordermodel.UpdateSupplyOrder, morekeys ...string) (*supplyordermodel.SupplyOrder, error) {

	record, err := biz.store.GetSupplyOrder(ctx, cond)
	if err != nil {
		return nil, common.ErrNotFoundEntity(supplyordermodel.EntityName, err)
	}

	if record == nil {
		return nil, common.ErrNotFoundEntity(supplyordermodel.EntityName, nil)
	}

	var updatedRecord *supplyordermodel.SupplyOrder
	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		var totalAmount float64 = 0
		newIngredientIDs := make(map[uint64]bool)
		for _, ing := range data.Ingredients {
			newIngredientIDs[ing.IngredientID] = true
			totalAmount += ing.Price * float64(ing.Quantity)
		}

		data.TotalAmount = totalAmount

		updatedRecord, err = biz.store.UpdateSupplyOrder(txCtx, map[string]interface{}{"supplyorder_id": record.SupplyOrderID}, data)

		if err != nil {
			return common.ErrCannotUpdateEntity(supplyordermodel.EntityName, err)
		}

		existingItemsMap := make(map[uint64]supplyordermodel.SupplyOrderItem)

		for _, item := range record.SupplyOrderItems {
			existingItemsMap[item.IngredientID] = item
		}

		// Danh sách các SupplyOrderItem cần cập nhật và tạo mới
		var itemsToUpdate []supplyordermodel.UpdateSupplyOrderItem
		var itemsToCreate []supplyordermodel.CreateSupplyOrderItem
		var stockBatchesToUpdate []stockbatchmodel.UpdateStockBatch
		var stockBatchesToCreate []stockbatchmodel.CreateStockBatch
		var stockBatchIDsToDelete []uint64

		for _, item := range data.Ingredients {
			if exitingItem, exits := existingItemsMap[item.IngredientID]; exits {
				supOrderItems := supplyordermodel.UpdateSupplyOrderItem{
					SupplyOrderItemID: exitingItem.SupplyOrderItemID,
					IngredientID:      item.IngredientID,
					Price:             item.Price,
					Quantity:          item.Quantity,
					Unit:              item.Unit,
					SupplyOrderID:     updatedRecord.SupplyOrderID,
					StockBatchID:      0,
				}

				itemsToUpdate = append(itemsToUpdate, supOrderItems)

				// Xử lý StockBatch
				stockBatchID, stockBatchUpdate, stockBatchCreate, err := biz.handleStockBatch(txCtx, item, exitingItem.StockBatchID)
				if err != nil {
					return common.ErrCannotUpdateEntity(supplyordermodel.EntityName, err)
				}

				supOrderItems.StockBatchID = uint64(stockBatchID)
				itemsToUpdate[len(itemsToUpdate)-1].StockBatchID = uint64(stockBatchID)

				if stockBatchUpdate != nil {
					stockBatchesToUpdate = append(stockBatchesToUpdate, *stockBatchUpdate)
				}
				if stockBatchCreate != nil {
					stockBatchesToCreate = append(stockBatchesToCreate, *stockBatchCreate)
				}
			} else {

				newItem := supplyordermodel.CreateSupplyOrderItem{
					IngredientID:  item.IngredientID,
					Price:         item.Price,
					Quantity:      item.Quantity,
					Unit:          item.Unit,
					SupplyOrderID: updatedRecord.SupplyOrderID,
					StockBatchID:  0, // Sẽ cập nhật sau khi tạo mới StockBatch
				}

				itemsToCreate = append(itemsToCreate, newItem)

				// Tạo mới StockBatch
				newStockBatch := stockbatchmodel.CreateStockBatch{
					Quantity:       item.Quantity,
					ExpirationDate: item.ExpirationDate,
					ReceivedDate:   item.ReceivedDate,
					IngredientID:   item.IngredientID,
				}
				stockBatchesToCreate = append(stockBatchesToCreate, newStockBatch)

			}

		}

		// Cập nhật các StockBatch đã thay đổi
		if len(stockBatchesToUpdate) > 0 {
			_, err := biz.storeStockBatch.UpdateStockBatches(txCtx, nil, stockBatchesToUpdate)
			if err != nil {
				return common.ErrCannotUpdateEntity("StockBatch", err)
			}
		}

		if len(stockBatchesToCreate) > 0 {
			createdStockBatchIDs, err := biz.storeStockBatch.CreateStockBatches(txCtx, stockBatchesToCreate)

			if err != nil {
				return common.ErrCannotUpdateEntity(stockbatchmodel.EntityName, err)
			}

			// Cập nhật StockBatchID cho các SupplyOrderItem đã tạo
			for i := range itemsToCreate {
				itemsToCreate[i].StockBatchID = createdStockBatchIDs[i]
			}

			// 2.7. Tạo mới các SupplyOrderItem
			if len(itemsToCreate) > 0 {
				err := biz.storeItem.CreateSupplyOrderItems(txCtx, itemsToCreate)
				if err != nil {
					return common.ErrCannotCreateEntity("SupplyOrderItem", err)
				}
			}
		}

		// 2.8. Cập nhật các SupplyOrderItem đã thay đổi
		if len(itemsToUpdate) > 0 {

			err := biz.storeItem.UpdateSupplyOrderItems(txCtx, itemsToUpdate)
			if err != nil {
				return common.ErrCannotUpdateEntity("SupplyOrderItem", err)
			}
		}

		var itemsToDelete []uint64
		for _, existingItem := range record.SupplyOrderItems {
			if !newIngredientIDs[existingItem.IngredientID] {
				itemsToDelete = append(itemsToDelete, existingItem.SupplyOrderItemID)
				stockBatchIDsToDelete = append(stockBatchIDsToDelete, existingItem.StockBatchID)
			}
		}

		if len(itemsToDelete) > 0 {
			// Xóa SupplyOrderItems
			err := biz.storeItem.DeleteSupplyOrderItems(txCtx, itemsToDelete)
			if err != nil {
				return common.ErrCannotDeleteEntity("SupplyOrderItem", err)
			}

			// Xóa StockBatch tương ứng
			err = biz.storeStockBatch.DeleteStockBatches(txCtx, stockBatchIDsToDelete)
			if err != nil {
				return common.ErrCannotDeleteEntity("StockBatch", err)
			}
		}

		{
			var pagging paggingcommon.Paging
			pagging.Process()

			// Generate cache key
			key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
				EntityName: supplyordermodel.EntityName,
				Cond:       cond,
				Paging:     pagging,
				Filter:     commonfilter.Filter{},
				MoreKeys:   morekeys,
			})
			if err != nil {
				common.ErrCannotGenerateKey(supplyordermodel.EntityName, err)
			}

			if err := biz.cache.DeleteSupplyOrder(ctx, key); err != nil {
				return common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
			}

			if err := biz.cache.DeleteListCache(txCtx, supplyordermodel.EntityName); err != nil {
				return common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
			}
		}

		return nil
	})

	return nil, nil
}

func (biz *updateSupplyOrderBusiness) handleStockBatch(ctx context.Context, ing supplyordermodel.CreateIngredient, existingStockBatchID uint64) (uint64, *stockbatchmodel.UpdateStockBatch, *stockbatchmodel.CreateStockBatch, error) {
	if existingStockBatchID != 0 {
		// Nếu StockBatch đã tồn tại, cập nhật Quantity, ExpirationDate, ReceivedDate
		stockBatch, err := biz.storeStockBatch.GetStockBatch(ctx, map[string]interface{}{"stockbatch_id": existingStockBatchID})
		if err != nil {
			return 0, nil, nil, common.ErrCannotGetEntity("StockBatch", err)
		}
		if stockBatch == nil {
			return 0, nil, nil, common.ErrNotFoundEntity("StockBatch", nil)
		}

		// Tạo struct UpdateStockBatch
		updatedQuantity := ing.Quantity
		update := stockbatchmodel.UpdateStockBatch{
			StockBatchID:   existingStockBatchID,
			Quantity:       updatedQuantity,
			ExpirationDate: ing.ExpirationDate,
			ReceivedDate:   ing.ReceivedDate,
		}

		return existingStockBatchID, &update, nil, nil
	} else {
		// Nếu StockBatch chưa tồn tại, tạo mới
		newStockBatch := stockbatchmodel.CreateStockBatch{
			Quantity:       ing.Quantity,
			ExpirationDate: ing.ExpirationDate,
			ReceivedDate:   ing.ReceivedDate,
			IngredientID:   ing.IngredientID,
		}

		return 0, nil, &newStockBatch, nil
	}
}
