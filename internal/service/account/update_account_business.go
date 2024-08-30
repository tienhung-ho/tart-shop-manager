package accountbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
	cacheutil "tart-shop-manager/internal/util/cache"
	responseutil "tart-shop-manager/internal/util/response"
)

type UpdateAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error)
}

type UpdateAccountCache interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	DeleteAccount(ctx context.Context, morekeys ...string) error
}

type updateAccountBusiness struct {
	store UpdateAccountStorage
	cache UpdateAccountCache
}

func NewUpdateAccount(store UpdateAccountStorage, cache UpdateAccountCache) *updateAccountBusiness {
	return &updateAccountBusiness{store, cache}
}

func (biz *updateAccountBusiness) UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error) {

	record, err := biz.store.GetAccount(ctx, cond)

	if err != nil {
		return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateAccount(ctx, map[string]interface{}{"account_id": record.AccountID}, account, morekeys...)

	if err != nil {
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, accountmodel.EntityName) // Extract field causing the duplicate error
			return nil, common.ErrDuplicateEntry(accountmodel.EntityName, fieldName, err)
		}

		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}
	var pagging paggingcommon.Paging
	pagging.Process()

	key := cacheutil.GenerateKey(accountmodel.EntityName, cond, pagging, commonfilter.Filter{})

	if err := biz.cache.DeleteAccount(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	return updatedRecord, nil
}
