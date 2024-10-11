package ingredientstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (s *mysqlIngredient) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter,
	moreKeys ...string) ([]ingredientmodel.Ingredient, error) {
	db := s.db.WithContext(ctx)

	// Count total records
	if err := s.countRecord(db, cond, paging, filter); err != nil {
		return nil, err
	}

	// Build query with filters
	query := s.buildQuery(db, cond, filter)

	// Add pagination and sorting
	query, err := s.addPaging(query, paging)
	if err != nil {
		return nil, err
	}

	// Execute query
	var records []ingredientmodel.Ingredient
	if err := query.Find(&records).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Trả về nil nếu không tìm thấy bản ghi nào
		}
		return nil, common.ErrDB(err)
	}

	return records, nil
}

func (s *mysqlIngredient) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	db = s.buildQuery(db, cond, filter)
	if err := db.Model(&ingredientmodel.Ingredient{}).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error counting items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlIngredient) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {
		if filter.Status != "" {
			db = db.Where("status IN ?", filter.Status)
		}
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
		}
		if len(filter.IDs) > 0 {
			db = db.Where("ingredient_id IN ?", filter.IDs)
		}
	}
	return db
}

func (s *mysqlIngredient) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
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
		db = db.Order("ingredient_id desc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}
