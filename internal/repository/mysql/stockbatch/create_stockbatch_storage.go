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

func (r *mysqlStockBatch) CreateStockBatch(ctx context.Context, data *stockbatchmodel.CreateStockBatch, morekeys ...string) (uint64, error) {

	db := r.db.Begin()

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

func (r *mysqlStockBatch) CreateStockBatches(ctx context.Context, data []stockbatchmodel.CreateStockBatch) ([]uint64, error) {
	if len(data) == 0 {
		return nil, nil
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, common.ErrDB(tx.Error)
	}

	defer commonrecover.RecoverTransaction(tx)

	// Thực hiện bulk insert
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				// Lỗi trùng lặp
				fieldName := responseutil.ExtractFieldFromError(err, stockbatchmodel.EntityName)
				tx.Rollback()
				return nil, common.ErrDuplicateEntry(stockbatchmodel.EntityName, fieldName, err)
			case 1452:
				// Lỗi khóa ngoại
				tx.Rollback()
				return nil, common.ErrForeignKeyConstraint(stockbatchmodel.EntityName, "ingredient_id", err)
			default:
				// Các lỗi MySQL khác
				tx.Rollback()
				return nil, common.ErrDB(err)
			}
		}
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Thu thập các StockBatchID đã được tạo
	stockIDs := make([]uint64, len(data))
	for i, batch := range data {
		stockIDs[i] = batch.StockBatchID
	}

	return stockIDs, nil
}
