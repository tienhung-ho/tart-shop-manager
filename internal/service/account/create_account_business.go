package accountbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
	responseutil "tart-shop-manager/internal/util/response"
)

type CreateAccountBusiness interface {
	CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error)
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
}

type createAccountBusiness struct {
	store CreateAccountBusiness
}

func NewCreateAccountbiz(store CreateAccountBusiness) *createAccountBusiness {
	return &createAccountBusiness{store: store}
}

func (biz *createAccountBusiness) CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error) {

	recordId, err := biz.store.CreateAccount(ctx, data, morekeys...)

	if err != nil {
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, accountmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(accountmodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	return recordId, nil
}
