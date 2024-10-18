package supplyorderstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (r *mysqlSupplyOrder) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]supplyordermodel.SupplyOrder, error) {
	db := r.db.WithContext(ctx)

	// Count total records
	if err := r.countRecord(db, cond, pagging, filter); err != nil {
		return nil, err
	}

	// Build query with filters
	query := r.buildQuery(db, cond, filter)

	// Add pagination and sorting
	query, err := r.addPaging(query, pagging)
	if err != nil {
		return nil, err
	}
	// Execute query
	var supplyOrders []supplyordermodel.SupplyOrder
	if err := query.
		Preload("SupplyOrderItems").
		Preload("SupplyOrderItems.StockBatch").
		Find(&supplyOrders).Error; err != nil {
		return nil, err
	}

	return supplyOrders, nil
}

func (r *mysqlSupplyOrder) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	// Apply conditions and filters
	db = r.buildQuery(db, cond, filter)

	if err := db.Model(&supplyordermodel.SupplyOrder{}).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error counting items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (r *mysqlSupplyOrder) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {
		if filter.Status != "" {
			db = db.Where("status IN ?", filter.Status)
		}

		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("description LIKE ?", searchPattern)
		}

		if filter.CategoryID != 0 {
			db = db.Where("supplyorder_id = ?", filter.CategoryID)
		}

		if filter.MinPrice > 0 {
			db = db.Where("total_amount >= ?", filter.MinPrice)
		}

		if filter.MaxPrice > 0 {
			db = db.Where("total_amount <= ?", filter.MaxPrice)
		}
	}
	return db
}

func (r *mysqlSupplyOrder) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
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
		db = db.Order("supplyorder_id desc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}
