package accountbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	DeleteAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteAccountCache interface {
	DeleteAccount(ctx context.Context, morekeys ...string) error
}

type deleteAccountBusiness struct {
	store DeleteAccountStorage
	cache DeleteAccountCache
}

func NewDeleteAccountBiz(store DeleteAccountStorage, cache DeleteAccountCache) *deleteAccountBusiness {
	return &deleteAccountBusiness{store, cache}
}

func (biz *deleteAccountBusiness) DeleteAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetAccount(ctx, cond)

	if err != nil {
		return common.ErrNotFoundEntity(accountmodel.EntityName, err)
	}

	if err := biz.store.DeleteAccount(ctx, map[string]interface{}{"account_id": record.AccountID}, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(accountmodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	key := cacheutil.GenerateKey(accountmodel.EntityName, cond, pagging, commonfilter.Filter{})

	if err := biz.cache.DeleteAccount(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(accountmodel.EntityName, err)
	}

	return nil
}
