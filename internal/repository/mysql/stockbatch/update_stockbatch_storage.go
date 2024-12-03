package stockbatchstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm/clause"
	"strings"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	responseutil "tart-shop-manager/internal/util/response"
)

// UpdateStockBatches updates multiple StockBatch records in a single query
// func (s *mysqlStockBatch) UpdateStockBatches(ctx context.Context, cond map[string]interface{},
//
//		data []stockbatchmodel.UpdateStockBatch) ([]uint64, error) {
//		if len(data) == 0 {
//			return nil, nil
//		}
//
//		// Khởi tạo transaction
//		tx := s.db.WithContext(ctx).Begin()
//		if tx.Error != nil {
//			return nil, common.ErrDB(tx.Error)
//		}
//
//		// Đảm bảo rollback nếu có lỗi
//		defer func() {
//			if r := recover(); r != nil {
//				tx.Rollback()
//			}
//		}()
//
//		var stockBatchIDs []uint64
//		var quantityCases []string
//		var expirationDateCases []string
//		var receivedDateCases []string
//		var updatedBy string
//
//		// Lấy email từ context
//		if email, ok := ctx.Value("email").(string); ok {
//			updatedBy = email
//		} else {
//			updatedBy = "system" // Hoặc giá trị mặc định khác
//		}
//
//		for _, d := range data {
//			stockBatchIDs = append(stockBatchIDs, d.StockBatchID)
//			quantityCases = append(quantityCases, fmt.Sprintf("WHEN %d THEN %f", d.StockBatchID, d.Quantity))
//			expirationDateFormatted := d.ExpirationDate.Format(common.DateFormat)
//			receivedDateFormatted := d.ReceivedDate.Format(common.DateFormat)
//			expirationDateCases = append(expirationDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, expirationDateFormatted))
//			receivedDateCases = append(receivedDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, receivedDateFormatted))
//		}
//
//		// Xây dựng các phần CASE cho từng trường
//		quantityCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(quantityCases, " "))
//		expirationDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(expirationDateCases, " "))
//		receivedDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(receivedDateCases, " "))
//
//		// Xây dựng WHERE clause
//		whereClause := "stockbatch_id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stockBatchIDs)), ","), "[]") + ")"
//
//		// Xây dựng câu lệnh SQL hoàn chỉnh, bao gồm cả việc cập nhật trường UpdatedBy
//		query := fmt.Sprintf(`
//	      UPDATE %s
//	      SET
//	          quantity = %s,
//	          expiration_date = %s,
//	          received_date = %s,
//	          updated_by = '%s'
//	      WHERE %s;
//	  `, stockbatchmodel.StockBatch{}.TableName(), quantityCaseStmt, expirationDateCaseStmt, receivedDateCaseStmt, updatedBy, whereClause)
//
//		// Thực thi câu lệnh SQL
//		if err := tx.Exec(query).Error; err != nil {
//			tx.Rollback()
//			return nil, common.ErrCannotUpdateEntity("StockBatch", err)
//		}
//
//		// Commit transaction nếu không có lỗi
//		if err := tx.Commit().Error; err != nil {
//			return nil, common.ErrDB(err)
//		}
//
//		return stockBatchIDs, nil
//	}
//
// UpdateStockBatches updates multiple StockBatch records in a single query
func (r *mysqlStockBatch) UpdateStockBatches(ctx context.Context, cond map[string]interface{},
	data []stockbatchmodel.UpdateStockBatch) ([]uint64, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Bắt đầu transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, common.ErrDB(tx.Error)
	}

	// Đảm bảo rollback nếu có lỗi
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var stockBatchIDs []uint64
	var quantityCases []string
	var expirationDateCases []string
	var receivedDateCases []string
	var updatedBy string

	// Lấy email từ context
	if email, ok := ctx.Value("email").(string); ok && email != "" {
		updatedBy = email
	} else {
		updatedBy = "system" // Hoặc giá trị mặc định khác
	}

	for _, d := range data {
		stockBatchIDs = append(stockBatchIDs, d.StockBatchID)
		if d.Quantity != nil {
			quantityCases = append(quantityCases, fmt.Sprintf("WHEN %d THEN %f", d.StockBatchID, *d.Quantity))
		}
		if d.ExpirationDate != nil {
			expirationDateFormatted := d.ExpirationDate.Format("2006-01-02") // Định dạng phù hợp
			expirationDateCases = append(expirationDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, expirationDateFormatted))
		}
		if d.ReceivedDate != nil {
			receivedDateFormatted := d.ReceivedDate.Format("2006-01-02") // Định dạng phù hợp
			receivedDateCases = append(receivedDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, receivedDateFormatted))
		}
	}

	// Xây dựng các phần CASE cho từng trường, chỉ bao gồm các trường không nil
	var setClauses []string
	//var args []interface{}

	if len(quantityCases) > 0 {
		quantityCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(quantityCases, " "))
		setClauses = append(setClauses, fmt.Sprintf("quantity = %s", quantityCaseStmt))
		// Thêm args cho quantity nếu cần
		// Do đã gán giá trị trực tiếp trong CASE, không cần thêm vào args
	}

	if len(expirationDateCases) > 0 {
		expirationDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(expirationDateCases, " "))
		setClauses = append(setClauses, fmt.Sprintf("expiration_date = %s", expirationDateCaseStmt))
		// Thêm args cho expiration_date nếu cần
	}

	if len(receivedDateCases) > 0 {
		receivedDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(receivedDateCases, " "))
		setClauses = append(setClauses, fmt.Sprintf("received_date = %s", receivedDateCaseStmt))
		// Thêm args cho received_date nếu cần
	}

	// Thêm cập nhật trường updated_by
	setClauses = append(setClauses, fmt.Sprintf("updated_by = '%s'", updatedBy))

	// Xây dựng WHERE clause an toàn
	var placeholders []string
	var whereArgs []interface{}
	for _, id := range stockBatchIDs {
		placeholders = append(placeholders, "?")
		whereArgs = append(whereArgs, id)
	}
	whereClause := fmt.Sprintf("stockbatch_id IN (%s)", strings.Join(placeholders, ","))

	// Xây dựng câu lệnh SQL hoàn chỉnh
	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf(`
       UPDATE %s
       SET
           %s
       WHERE %s;
   `, stockbatchmodel.UpdateStockBatch{}.TableName(), setClause, whereClause)

	// Thêm log để kiểm tra câu lệnh SQL và các args
	fmt.Printf("Executing SQL: %s with args: %v\n", query, whereArgs)

	// Thực thi câu lệnh SQL với args
	if err := tx.Exec(query, whereArgs...).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrCannotUpdateEntity("StockBatch", err)
	}

	// Commit transaction nếu không có lỗi
	if err := tx.Commit().Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return stockBatchIDs, nil
}

func (r *mysqlStockBatch) UpdateStockBatch(ctx context.Context, cond map[string]interface{},
	data *stockbatchmodel.UpdateStockBatch) (*stockbatchmodel.StockBatch, error) {

	db := r.db.Begin()

	if db.Error != nil {
		return nil, common.ErrDB(db.Error)
	}

	defer commonrecover.RecoverTransaction(db)

	if err := db.WithContext(ctx).Model(&stockbatchmodel.UpdateStockBatch{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, stockbatchmodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(stockbatchmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return nil, err
	}

	var record stockbatchmodel.StockBatch

	if err := db.WithContext(ctx).Model(data).
		Where(cond).
		Preload("Ingredient").
		First(&record).Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &record, nil
}
