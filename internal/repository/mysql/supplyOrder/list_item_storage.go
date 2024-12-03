package supplyorderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

//func (r *mysqlSupplyOrder) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error) {
//	db := r.db.WithContext(ctx)
//
//	// Count total records
//	if err := r.countRecord(db, cond, pagging, filter); err != nil {
//		return nil, err
//	}
//
//	// Build query with filters
//	query := r.buildQuery(db, cond, filter)
//
//	// Add pagination and sorting
//	query, err := r.addPaging(query, pagging)
//	if err != nil {
//		return nil, err
//	}
//	// Execute query
//	var supplyOrders []supplyordermodel.SupplyOrder
//	if err := query.
//		Preload("SupplyOrderItems").
//		Preload("SupplyOrderItems.StockBatch").
//		Find(&supplyOrders).Error; err != nil {
//		return nil, err
//	}
//
//	return supplyOrders, nil
//}
//
//func (r *mysqlSupplyOrder) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
//	// Apply conditions and filters
//	db = r.buildQuery(db, cond, filter)
//
//	if err := db.Model(&supplyordermodel.SupplyOrder{}).Count(&paging.Total).Error; err != nil {
//		return common.NewErrorResponse(err, "Error counting items from database", err.Error(), "CouldNotCount")
//	}
//	return nil
//}
//
//func (r *mysqlSupplyOrder) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
//	db = db.Where(cond)
//	if filter != nil {
//		if filter.Status != "" {
//			db = db.Where("status IN ?", filter.Status)
//		}
//
//		if filter.Search != "" {
//			searchPattern := "%" + filter.Search + "%"
//			db = db.Where("description LIKE ?", searchPattern)
//		}
//
//		if filter.CategoryID != 0 {
//			db = db.Where("supplyorder_id = ?", filter.CategoryID)
//		}
//
//		if filter.MinPrice > 0 {
//			db = db.Where("total_amount >= ?", filter.MinPrice)
//		}
//
//		if filter.MaxPrice > 0 {
//			db = db.Where("total_amount <= ?", filter.MaxPrice)
//		}
//	}
//	return db
//}
//
//func (r *mysqlSupplyOrder) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
//	// Parse and validate the sort fields
//	sortFields, err := paging.ParseSortFields(paging.Sort, AllowedSortFields)
//	if err != nil {
//		return nil, common.NewErrorResponse(err, "Invalid sort parameters", err.Error(), "InvalidSort")
//	}
//
//	// Apply sorting to the query
//	if len(sortFields) > 0 {
//		for _, sortField := range sortFields {
//			db = db.Order(sortField)
//		}
//	} else {
//		// Default sorting if no sort parameters are provided
//		db = db.Order("supplyorder_id desc")
//	}
//
//	// Apply pagination
//	offset := (paging.Page - 1) * paging.Limit
//	db = db.Offset(offset).Limit(paging.Limit)
//
//	return db, nil
//}

func (r *mysqlSupplyOrder) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error) {

	db := r.db.WithContext(ctx)

	// Build base query
	query := r.buildBaseQuery(db, cond)

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total records
	if err := r.countTotalRecords(query, paging); err != nil {
		return nil, err
	}

	// Execute main query with preload
	supplyOrders, err := r.executeMainQuery(query, paging, filter)
	if err != nil {
		return nil, err
	}

	return supplyOrders, nil
}

func (r *mysqlSupplyOrder) buildBaseQuery(db *gorm.DB, cond map[string]interface{}) *gorm.DB {
	return db.Model(&supplyordermodel.SupplyOrder{}).
		Where(cond)
}

func (r *mysqlSupplyOrder) applyFilters(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.IngredientID != nil {
		query = query.Where("ingredient_id = ?", filter.IngredientID)
	}

	if filter.Search != "" {
		query = query.Where("description LIKE ?", "%"+filter.Search+"%")
	}

	if filter.MinPrice > 0 {
		query = query.Where("total_amount >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("total_amount <= ?", filter.MaxPrice)
	}

	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("order_date BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	} else if filter.StartDate != nil {
		query = query.Where("order_date >= ?", filter.StartDate)
	} else if filter.EndDate != nil {
		query = query.Where("order_date <= ?", filter.EndDate)
	}

	// Apply date filters on associated StockBatch
	if filter.StartExpirationDate != nil || filter.EndExpirationDate != nil ||
		filter.StartReceivedDate != nil || filter.EndReceivedDate != nil ||
		filter.ExpirationDate != nil || filter.ReceivedDate != nil {

		subQuery := r.db.Model(&supplyordermodel.SupplyOrderItem{}).
			Select("supplyorder_id").
			Joins("JOIN StockBatch ON StockBatch.stockbatch_id = SupplyOrderItem.stockbatch_id").
			Scopes(func(db *gorm.DB) *gorm.DB {
				return r.applyDateFilter(db, filter)
			})

		query = query.Where("SupplyOrder.supplyorder_id IN (?)", subQuery)
	}

	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlSupplyOrder) applyDateFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter.StartExpirationDate != nil {
		query = query.Where("StockBatch.expiration_date >= ?", filter.StartExpirationDate)
	}

	if filter.EndExpirationDate != nil {
		query = query.Where("StockBatch.expiration_date <= ?", filter.EndExpirationDate)
	}

	if filter.StartReceivedDate != nil {
		query = query.Where("StockBatch.received_date >= ?", filter.StartReceivedDate)
	}

	if filter.EndReceivedDate != nil {
		query = query.Where("StockBatch.received_date <= ?", filter.EndReceivedDate)
	}

	if filter.ExpirationDate != nil {
		query = query.Where("StockBatch.expiration_date = ?", filter.ExpirationDate)
	}

	if filter.ReceivedDate != nil {
		query = query.Where("StockBatch.received_date = ?", filter.ReceivedDate)
	}

	return query
}

func (r *mysqlSupplyOrder) applyOrderDateFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	//if filter.StartExpirationDate != nil {
	//	query = query.Where("(order_date = ?)", filter.OrderDate)
	//}

	return query
}

// countTotalRecords đếm tổng số records
func (r *mysqlSupplyOrder) countTotalRecords(query *gorm.DB, paging *paggingcommon.Paging) error {
	var total int64
	if err := query.Model(&supplyordermodel.SupplyOrder{}).Count(&total).Error; err != nil {
		return common.ErrDB(err)
	}
	paging.Total = total
	return nil
}

// executeMainQuery thực hiện query chính với preload
func (r *mysqlSupplyOrder) executeMainQuery(
	query *gorm.DB,
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter,
) ([]supplyordermodel.SupplyOrder, error) {
	var supplyOrders []supplyordermodel.SupplyOrder

	err := query.
		Preload("SupplyOrderItems", r.supplyOrderItemsPreloadCondition(filter)).
		Preload("SupplyOrderItems.StockBatch", r.stockBatchPreloadCondition(filter)).
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&supplyOrders).Error

	if err != nil {
		return nil, common.ErrDB(err)
	}

	return supplyOrders, nil
}

// supplyOrderItemsPreloadCondition tạo điều kiện preload cho SupplyOrderItems
func (r *mysqlSupplyOrder) supplyOrderItemsPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyOrderDateFilter(db, filter)
	}
}

// stockBatchPreloadCondition tạo điều kiện preload cho StockBatch
func (r *mysqlSupplyOrder) stockBatchPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyDateFilter(db, filter)
	}
}
