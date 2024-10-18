package stockbatchstorage

import (
	"context"
	"fmt"
	"strings"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

// UpdateStockBatches updates multiple StockBatch records in a single query
func (s *mysqlStockBatch) UpdateStockBatches(ctx context.Context, cond map[string]interface{}, data []stockbatchmodel.UpdateStockBatch) ([]uint64, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Khởi tạo transaction
	tx := s.db.WithContext(ctx).Begin()
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
	if email, ok := ctx.Value("email").(string); ok {
		updatedBy = email
	} else {
		updatedBy = "system" // Hoặc giá trị mặc định khác
	}

	for _, d := range data {
		stockBatchIDs = append(stockBatchIDs, d.StockBatchID)
		quantityCases = append(quantityCases, fmt.Sprintf("WHEN %d THEN %d", d.StockBatchID, d.Quantity))
		expirationDateFormatted := d.ExpirationDate.Format(common.DateFormat)
		receivedDateFormatted := d.ReceivedDate.Format(common.DateFormat)
		expirationDateCases = append(expirationDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, expirationDateFormatted))
		receivedDateCases = append(receivedDateCases, fmt.Sprintf("WHEN %d THEN '%s'", d.StockBatchID, receivedDateFormatted))
	}

	// Xây dựng các phần CASE cho từng trường
	quantityCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(quantityCases, " "))
	expirationDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(expirationDateCases, " "))
	receivedDateCaseStmt := fmt.Sprintf("CASE stockbatch_id %s END", strings.Join(receivedDateCases, " "))

	// Xây dựng WHERE clause
	whereClause := "stockbatch_id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stockBatchIDs)), ","), "[]") + ")"

	// Xây dựng câu lệnh SQL hoàn chỉnh, bao gồm cả việc cập nhật trường UpdatedBy
	query := fmt.Sprintf(`
        UPDATE %s
        SET 
            quantity = %s,
            expiration_date = %s,
            received_date = %s,
            updated_by = '%s'
        WHERE %s;
    `, stockbatchmodel.StockBatch{}.TableName(), quantityCaseStmt, expirationDateCaseStmt, receivedDateCaseStmt, updatedBy, whereClause)

	// Thực thi câu lệnh SQL
	if err := tx.Exec(query).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrCannotUpdateEntity("StockBatch", err)
	}

	// Commit transaction nếu không có lỗi
	if err := tx.Commit().Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return stockBatchIDs, nil
}
