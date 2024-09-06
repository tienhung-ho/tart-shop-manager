package permissionstorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	commonrecover "tart-shop-manager/internal/common/recover"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

func (s *mysqlPermission) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]permissionmodel.Permission, error) {

	db := s.db

	defer commonrecover.RecoverTransaction(db)

	// Đếm tổng số lượng items
	if err := s.countRecord(db, cond, paging, filter); err != nil {
		return nil, err
	}

	// Xây dựng truy vấn động
	query := s.buildQuery(db, cond, filter)

	// Thêm phân trang
	query = s.addPaging(query, paging)

	// Thực hiện truy vấn
	var records []permissionmodel.Permission
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *mysqlPermission) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	if permissionIds, is := cond["permission_id"]; is {

		if permissionList, valid := permissionIds.([]uint); valid && len(permissionList) > 0 {
			db = db.Where("permission_id IN ?", permissionList)
		} else {
			return common.NewErrorResponse(nil, "Invalid names format", "names must be a non-empty slice of strings", "InvalidNames")
		}

	} else {
		db = db.Where(cond)
		if filter != nil && filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}
	}
	if err := db.Table(permissionmodel.Permission{}.TableName()).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error count items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlPermission) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil && filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	return db
}

func (s *mysqlPermission) addPaging(db *gorm.DB, paging *paggingcommon.Paging) *gorm.DB {
	return db.Order("permission_id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit)
}
