package reportstorage

import (
	"context"
	"gorm.io/gorm"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	reportmodel "tart-shop-manager/internal/entity/dtos/sql/report"
	"time"
)

func (r *mysqlReportOrder) ReportOrder(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]reportmodel.Order, error) {

	db := r.db.WithContext(ctx)

	// Build base query
	query := r.buildBaseQuery(db, cond)

	// Apply filters
	query = r.applyFilters(query, filter)

	query, err := r.addPaging(query, paging)
	if err != nil {
		return nil, err
	}

	// Count total records
	if err := r.countTotalRecords(query, paging); err != nil {
		return nil, err
	}

	// Execute main query with preload
	orders, err := r.executeMainQuery(query, paging, filter)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *mysqlReportOrder) buildBaseQuery(db *gorm.DB, cond map[string]interface{}) *gorm.DB {
	return db.Table("`Order`").
		Where(cond)
}

func (r *mysqlReportOrder) applyFilters(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.MinPrice > 0 {
		query = query.Where("`Order`.total_amount >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("`Order`.total_amount <= ?", filter.MaxPrice)
	}

	loc, err := time.LoadLocation("UTC") // Hoặc múi giờ bạn sử dụng
	if err != nil {
		log.Print("Error loading location:", err)
		return query
	}

	if filter.InDate != nil {
		startOfDay := time.Date(filter.InDate.Year(), filter.InDate.Month(), filter.InDate.Day(), 0, 0, 0, 0, loc)
		endOfDay := startOfDay.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		query = query.Where("`Order`.order_date >= ? ", filter.InDate)
		query = query.Where("`Order`.order_date < ? ", endOfDay)
		log.Print(startOfDay, endOfDay)
	}

	if filter.StartDate != nil {
		query = query.Where("`Order`.order_date >= ?", filter.StartDate)
	}

	if filter.EndDate != nil {
		endOfDay := filter.EndDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		query = query.Where("`Order`.order_date <= ?", endOfDay)
	}

	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlReportOrder) applyOrderItemFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlReportOrder) applyRecipeFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	return query
}

func (r *mysqlReportOrder) applyProductFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	return query
}

// countTotalRecords đếm tổng số records
func (r *mysqlReportOrder) countTotalRecords(query *gorm.DB, paging *paggingcommon.Paging) error {
	var total int64
	if err := query.Model(&reportmodel.Order{}).Count(&total).Error; err != nil {
		return common.ErrDB(err)
	}
	paging.Total = total
	return nil
}

// executeMainQuery thực hiện query chính với preload
func (r *mysqlReportOrder) executeMainQuery(
	query *gorm.DB,
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter,
) ([]reportmodel.Order, error) {
	var orders []reportmodel.Order

	err := query.
		Preload("OrderItems", r.orderItemsPreloadCondition(filter)).
		Preload("OrderItems.Recipe", r.recipePreloadCondition(filter)).
		Preload("OrderItems.Recipe.Product", r.orderItemsRecipeProductCondition(filter)).
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&orders).Error

	if err != nil {
		return nil, common.ErrDB(err)
	}

	return orders, nil
}

// supplyOrderItemsPreloadCondition tạo điều kiện preload cho SupplyOrderItems
func (r *mysqlReportOrder) orderItemsPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyOrderItemFilter(db, filter)
	}
}

func (r *mysqlReportOrder) orderItemsRecipeProductCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyProductFilter(db, filter)
	}
}

// stockBatchPreloadCondition tạo điều kiện preload cho StockBatch
func (r *mysqlReportOrder) recipePreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyRecipeFilter(db, filter)
	}
}

func (s *mysqlReportOrder) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
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
