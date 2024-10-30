package orderitemstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlOrderItem) CreateOrderItems(ctx context.Context, data []ordermodel.CreateOrderItem) error {

	if len(data) == 0 {
		return nil
	}

	tx := s.getDB(ctx)
	if tx.Error != nil {
		return common.ErrDB(tx.Error)
	}

	defer commonrecover.RecoverTransaction(tx)

	// Thực hiện bulk insert
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				// Lỗi trùng lặp
				fieldName := responseutil.ExtractFieldFromError(err, ordermodel.EntityNameOrderItem)
				return common.ErrDuplicateEntry(ordermodel.EntityNameOrderItem, fieldName, err)
			case 1452:
				// Lỗi khóa ngoại
				return common.ErrForeignKeyConstraint(ordermodel.EntityNameOrderItem, "recipe_id", err) // Sửa "ingredient_id" thành "recipe_id"
			default:
				// Các lỗi MySQL khác
				return common.ErrDB(err)
			}
		}
		return common.ErrDB(err)
	}

	// Không gọi tx.Commit() hay tx.Rollback() ở đây

	return nil
}
