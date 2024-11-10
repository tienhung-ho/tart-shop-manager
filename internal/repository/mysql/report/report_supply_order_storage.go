package reportstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	reportmodel "tart-shop-manager/internal/entity/dtos/sql/report"
	"time"
)

func (r *mysqlReportOrder) ReportSupplyOrder(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]reportmodel.SupplyOrder, error) {

	db := r.db.WithContext(ctx)

	// Build base query
	query := r.buildBaseQuerySupply(db, cond)

	// Apply filters
	query = r.applyFiltersSupply(query, filter)

	query, err := r.addPagingSupply(query, paging)
	if err != nil {
		return nil, err
	}

	// Count total records
	if err := r.countTotalRecordsSupply(query, paging); err != nil {
		return nil, err
	}

	// Execute main query with preload
	supplyOrders, err := r.executeMainQuerySupply(query, paging, filter)
	if err != nil {
		return nil, err
	}

	return supplyOrders, nil
}

func (r *mysqlReportOrder) buildBaseQuerySupply(db *gorm.DB, cond map[string]interface{}) *gorm.DB {
	return db.Table("SupplyOrder").
		Where(cond)
}

func (r *mysqlReportOrder) applyFiltersSupply(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.Ingredient != nil {
		query = query.Where("ingredient_id = ?", filter.Ingredient)
	}

	if filter.MinPrice > 0 {
		query = query.Where("SupplyOrder.total_amount >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("SupplyOrder.total_amount <= ?", filter.MaxPrice)
	}

	if filter.InDate != nil {
		query = query.Where("`SupplyOrder`.order_date >= ? ", filter.InDate)
		endOfDay := filter.InDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		query = query.Where("`SupplyOrder`.order_date <= ? ", endOfDay)
	}

	if filter.StartDate != nil {
		query = query.Where("`SupplyOrder`.order_date >= ?", filter.StartDate)
	}

	if filter.EndDate != nil {
		endOfDay := filter.EndDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		query = query.Where("`SupplyOrder`.order_date <= ?", endOfDay)
	}

	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlReportOrder) applyDateFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {

	return query
}

// countTotalRecordsSupply đếm tổng số records
func (r *mysqlReportOrder) countTotalRecordsSupply(query *gorm.DB, paging *paggingcommon.Paging) error {
	var total int64
	if err := query.Model(&reportmodel.SupplyOrder{}).Count(&total).Error; err != nil {
		return common.ErrDB(err)
	}
	paging.Total = total
	return nil
}

// executeMainQuerySupply thực hiện query chính với preload
func (r *mysqlReportOrder) executeMainQuerySupply(
	query *gorm.DB,
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter,
) ([]reportmodel.SupplyOrder, error) {
	var supplyOrders []reportmodel.SupplyOrder

	err := query.
		Preload("SupplyOrderItems", r.supplyOrderItemsPreloadCondition(filter)).
		Preload("SupplyOrderItems.Ingredient", r.ingredientPreloadCondition(filter)).
		Preload("Supplier", r.ingredientPreloadCondition(filter)).
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&supplyOrders).Error

	if err != nil {
		return nil, common.ErrDB(err)
	}

	return supplyOrders, nil
}

// supplyOrderItemsPreloadCondition tạo điều kiện preload cho SupplyOrderItems
func (r *mysqlReportOrder) supplyOrderItemsPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyDateFilter(db, filter)
	}
}

func (r *mysqlReportOrder) ingredientPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyDateFilter(db, filter)
	}
}

func (r *mysqlReportOrder) supplierPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyDateFilter(db, filter)
	}
}

func (s *mysqlReportOrder) addPagingSupply(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
	// Parse and validate the sort fields
	sortFields, err := paging.ParseSortFields(paging.Sort, AllowedSortFields)
	if err != nil {
		return nil, common.NewErrorResponse(err, "Invalid sort parameters", err.Error(), "InvalidSort")
	}

	// Apply sorting to the query
	if len(sortFields) > 0 {
		for _, sortField := range sortFields {
			db = db.Order(sortField)
		}
	} else {
		// Default sorting if no sort parameters are provided
		db = db.Order("order_date desc")
	}
	return db, nil
}
