package orderstorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlOrder) CreateOrder(ctx context.Context, data *ordermodel.CreateOrder) (uint64, error) {

	db := s.getDB(ctx)

	email, ok := ctx.Value("email").(string)

	if !ok {
		data.UpdatedBy = "system" // Hoặc giá trị mặc định khác
	}

	data.UpdatedBy = email

	accountID, ok := ctx.Value("id").(uint64)

	if !ok {
		data.AccountID = 1 // Hoặc giá trị mặc định khác
	}

	data.AccountID = accountID

	if err := db.Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := responseutil.ExtractFieldFromError(err, ordermodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(ordermodel.EntityName, fieldName, err)
		}
		return 0, common.ErrDB(err)
	}

	return data.OrderID, nil
}
