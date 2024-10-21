package stockbatchstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (r *mysqlStockBatch) ListItem(ctx context.Context, cond map[string]interface{},
	pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]stockbatchmodel.StockBatch, error) {
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
	var stockBatches []stockbatchmodel.StockBatch
	if err := query.
		Preload("Ingredient").
		Find(&stockBatches).Error; err != nil {
		return nil, err
	}

	return stockBatches, nil
}

func (r *mysqlStockBatch) countRecord(db *gorm.DB, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	// Apply conditions and filters
	db = r.buildQuery(db, cond, filter)

	if err := db.Model(&stockbatchmodel.StockBatch{}).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error counting items from database",
			err.Error(), "CouldNotCount")
	}
	return nil
}

func (r *mysqlStockBatch) buildQuery(db *gorm.DB, cond map[string]interface{},
	filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {
		if filter.Status != "" {
			db = db.Where("status IN ?", filter.Status)
		}

		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
		}

		if filter.Ingredient != nil {
			db = db.Where("ingredient_id = ?", filter.Ingredient)
		}

		// Thêm điều kiện lọc theo khoảng thời gian
		if filter.StartDate != nil && filter.EndDate != nil {
			// Giả sử bạn muốn lọc theo trường `received_date`
			db = db.Where("expiration_date BETWEEN ? AND ?", *filter.StartDate, *filter.EndDate)
		} else if filter.StartDate != nil {
			db = db.Where("expiration_date >= ?", *filter.StartDate)
		} else if filter.EndDate != nil {
			db = db.Where("expiration_date <= ?", *filter.EndDate)
		}
	}
	return db
}

func (r *mysqlStockBatch) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
	// Parse and validate the sort fields
	sortFields, err := paging.ParseSortFields(paging.Sort, AllowedSortFields)
	if err != nil {
		return nil, common.NewErrorResponse(err,
			"Invalid sort parameters", err.Error(), "InvalidSort")
	}

	// Apply sorting to the query
	if len(sortFields) > 0 {
		for _, sortField := range sortFields {
			db = db.Order(sortField)
		}
	} else {
		// Default sorting if no sort parameters are provided
		db = db.Order("stockbatch_id desc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}
