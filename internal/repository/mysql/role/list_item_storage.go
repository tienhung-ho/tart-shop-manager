package rolestorage

import (
	"context"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	commonrecover "tart-shop-manager/internal/common/recover"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

func (s *mysqlRole) ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error) {

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
	var records []rolemodel.Role
	if err := query.
		Select(rolemodel.SelectFields).
		Preload("Permissions").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *mysqlRole) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter) error {
	if roleIdS, is := cond["role_id"]; is {

		if roleList, valid := roleIdS.([]uint); valid && len(roleList) > 0 {
			db = db.Where("role_id IN ?", roleList)
		} else {
			return common.NewErrorResponse(nil, "Invalid names format", "names must be a non-empty slice of strings", "InvalidNames")
		}

	} else {
		db = db.Where(cond)
		if filter != nil && filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}
	}
	if err := db.Table(rolemodel.Role{}.TableName()).Count(&paging.Total).Error; err != nil {
		return common.NewErrorResponse(err, "Error count items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlRole) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *commonfilter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil && filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	return db
}

func (s *mysqlRole) addPaging(db *gorm.DB, paging *paggingcommon.Paging) *gorm.DB {
	return db.Order("role_id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit)
}
