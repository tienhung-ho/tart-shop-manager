package stockbatchstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlStockBatch) CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint, error) {

	db := s.db.Begin()

	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				fieldName := responseutil.ExtractFieldFromError(err, stockbatchmodel.EntityName)
				return 0, common.ErrDuplicateEntry(stockbatchmodel.EntityName, fieldName, err)
			case 1452:
				// Xử lý lỗi khóa ngoại
				return 0, common.ErrForeignKeyConstraint(stockbatchmodel.EntityName, "ingredient_id", err)
			default:
				// Các lỗi MySQL khác
				db.Rollback()
				return 0, err
			}
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.StockBatchID, nil
}
