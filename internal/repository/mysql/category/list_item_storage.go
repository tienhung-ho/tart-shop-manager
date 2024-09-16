package categorystorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	commonrecover "tart-shop-manager/internal/common/recover"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

func (s *mysqlCategory) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]categorymodel.Category, error) {
	db := s.db

	defer commonrecover.RecoverTransaction(db)

	// // Đếm tổng số lượng items
	if err := s.countRecord(db, cond, paging, filter); err != nil {
		return nil, err
	}

	// Xây dựng truy vấn động
	query := s.buildQuery(db, cond, filter)

	// Thêm phân trang
	query = s.addPaging(query, paging)

	// Thực hiện truy vấn
	var records []categorymodel.Category
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *mysqlCategory) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	if names, ok := cond["names"]; ok {
		db = db.Where("name IN ?", names)
	} else {
		db = db.Where(cond)
		if filter != nil && filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}
	}
	if err := db.Table(categorymodel.Category{}.TableName()).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error count items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlCategory) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil && filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	return db
}

func (s *mysqlCategory) addPaging(db *gorm.DB, paging *paggingcommon.Paging) *gorm.DB {
	return db.Order("category_id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit)
}
