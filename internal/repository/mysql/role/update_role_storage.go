package rolestorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlRole) UpdateRole(ctx context.Context, cond map[string]interface{}, data *rolemodel.UpdateRole, morekeys ...string) error {

	// Bắt đầu transaction
	db := s.db.Begin()
	if db.Error != nil {
		return common.ErrDB(db.Error)
	}

	// Defer để đảm bảo rollback nếu có lỗi
	defer commonrecover.RecoverTransaction(db)

	// Cập nhật thông tin role
	if err := db.WithContext(ctx).Model(&rolemodel.UpdateRole{}).Where(cond).Updates(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, ingredientmodel.EntityName) // Extract field causing the duplicate error
			return common.ErrDuplicateEntry(ingredientmodel.EntityName, fieldName, err)
		}
		db.Rollback() // Rollback ngay khi có lỗi
		return err
	}

	// Nếu có permissions mới, xử lý cập nhật bảng nối
	//log.Print(data.Permissions != nil)
	//if data.Permissions != nil {
	var role rolemodel.Role

	// Lấy thông tin role
	if err := db.WithContext(ctx).Where(cond).First(&role).Error; err != nil {
		db.Rollback() // Rollback nếu không tìm thấy role
		return common.ErrDB(err)
	}

	// Thay thế các permissions mới
	if err := db.Model(&role).Association("Permissions").Replace(data.Permissions); err != nil {
		db.Rollback() // Rollback nếu thay thế thất bại
		return common.ErrDB(err)
	}
	//}

	// Commit transaction nếu không có lỗi
	if err := db.Commit().Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
