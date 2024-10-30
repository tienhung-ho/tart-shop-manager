package recipestorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

func (s *mysqlRecipe) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]recipemodel.Recipe, error) {
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
	var records []recipemodel.Recipe
	if err := query.
		Preload("Product").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").
		Find(&records).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Xử lý khi không tìm thấy bản ghi
			return nil, err
		}
		return nil, common.ErrDB(err)
	}

	return records, nil
}

func (s *mysqlRecipe) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	// Apply conditions and filters
	db = s.buildQuery(db, cond, filter)

	if err := db.Model(&recipemodel.Recipe{}).
		Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error counting items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlRecipe) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {
		if filter.Status != "" {
			db = db.Where("status IN ?", filter.Status)
		}

		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
		}

		if filter.CategoryID != 0 {
			db = db.Where("category_id = ?", filter.CategoryID)
		}

		if filter.MinPrice > 0 {
			db = db.Where("price >= ?", filter.MinPrice)
		}

		if filter.MaxPrice > 0 {
			db = db.Where("price <= ?", filter.MaxPrice)
		}

		if len(filter.IDs) > 0 {
			db = db.Where("recipe_id IN (?)", filter.IDs)
		}

		if len(filter.ProductIDs) > 0 {
			db = db.Where("product_id IN (?)", filter.ProductIDs)
		}

		// Thêm điều kiện lọc theo danh sách size
		if len(filter.Sizes) > 0 {
			db = db.Where("size IN (?)", filter.Sizes)
		}
	}
	return db
}

func (s *mysqlRecipe) addPaging(db *gorm.DB, paging *paggingcommon.Paging) (*gorm.DB, error) {
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
		db = db.Order("recipe_id desc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}
