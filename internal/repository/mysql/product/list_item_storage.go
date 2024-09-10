package productstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	commonrecover "tart-shop-manager/internal/common/recover"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (s *mysqlProduct) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error) {
	db := s.db

	defer commonrecover.RecoverTransaction(db)

	// Đếm tổng số lượng items
	if err := s.countRecord(db, cond, pagging, filter); err != nil {
		return nil, err
	}

	// Xây dựng truy vấn động
	query := s.buildQuery(db, cond, filter)

	// Thêm phân trang
	query = s.addPaging(query, pagging)

	// Thực hiện truy vấn
	var records []productmodel.Product
	if err := query.Preload("Category").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *mysqlProduct) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	if productIds, is := cond["product_id"]; is {

		if productList, valid := productIds.([]uint); valid && len(productList) > 0 {
			db = db.Where("product_id IN ?", productList)
		} else {
			return common.NewErrorResponse(nil, "Invalid names format", "names must be a non-empty slice of strings", "InvalidNames")
		}

	} else {
		db = db.Where(cond)
		if filter != nil && filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}
	}
	if err := db.Table(productmodel.Product{}.TableName()).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error count items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlProduct) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil && filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	return db
}

func (s *mysqlProduct) addPaging(db *gorm.DB, paging *paggingcommon.Paging) *gorm.DB {
	return db.Order("product_id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit)
}
