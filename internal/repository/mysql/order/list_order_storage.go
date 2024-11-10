package orderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

func (r *mysqlOrder) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) ([]ordermodel.Order, error) {

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

func (r *mysqlOrder) buildBaseQuery(db *gorm.DB, cond map[string]interface{}) *gorm.DB {
	return db.Table("`Order`").
		Where(cond)
}

func (r *mysqlOrder) applyFilters(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.MinPrice > 0 {
		query = query.Where("`Order`.total_amount >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("`Order`.total_amount <= ?", filter.MaxPrice)
	}

	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlOrder) applyOrderItemFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter.MinPrice != -1 {
		query = query.Where("(`OrderRecipe`.price >= ?)", filter.MinPrice)
	}

	if filter.MaxPrice != 0 {
		query = query.Where("(`OrderRecipe`.price <= ?)", filter.MaxPrice)
	}

	return query
}

// applyDateFilter áp dụng bộ lọc ngày
func (r *mysqlOrder) applyRecipeFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter.Recipe.Sizes != nil {
		query = query.Where("(Recipe.size = ?)", filter.Recipe.Sizes)
	}

	if filter.Recipe.ProductIDs != nil {
		query = query.Where("product_id IN ?", filter.Recipe.ProductIDs)
	}

	return query
}

func (r *mysqlOrder) applyProductFilter(query *gorm.DB, filter *commonfilter.Filter) *gorm.DB {
	if filter.Product.Name != "" {
		query = query.Where("(Product.name LIKE ?)", filter.Product.Name)
	}

	return query
}

// countTotalRecords đếm tổng số records
func (r *mysqlOrder) countTotalRecords(query *gorm.DB, paging *paggingcommon.Paging) error {
	var total int64
	if err := query.Model(&ordermodel.Order{}).Count(&total).Error; err != nil {
		return common.ErrDB(err)
	}
	paging.Total = total
	return nil
}

func (s *mysqlOrder) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
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
		db = db.Order("order_id desc")
	}
	return db, nil
}

// executeMainQuery thực hiện query chính với preload
func (r *mysqlOrder) executeMainQuery(
	query *gorm.DB,
	paging *paggingcommon.Paging,
	filter *commonfilter.Filter,
) ([]ordermodel.Order, error) {
	var orders []ordermodel.Order

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
func (r *mysqlOrder) orderItemsPreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyOrderItemFilter(db, filter)
	}
}

func (r *mysqlOrder) orderItemsRecipeProductCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyProductFilter(db, filter)
	}
}

// stockBatchPreloadCondition tạo điều kiện preload cho StockBatch
func (r *mysqlOrder) recipePreloadCondition(filter *commonfilter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		return r.applyRecipeFilter(db, filter)
	}
}
