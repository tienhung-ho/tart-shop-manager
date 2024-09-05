package accountbusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	"tart-shop-manager/internal/entity/dtos/sql/account"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
}

type GetAccountRedisStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	SaveAccount(ctx context.Context, data interface{}, morekeys ...string) error
}

type getAccountBusiness struct {
	store    GetAccountStorage
	rdbStore GetAccountRedisStorage
}

func NewGetAccountBiz(store GetAccountStorage, rdbStore GetAccountRedisStorage) *getAccountBusiness {
	return &getAccountBusiness{
		store:    store,
		rdbStore: rdbStore,
	}
}

func (biz *getAccountBusiness) GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error) {

	record, err := biz.rdbStore.GetAccount(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrCannotGetEntity(accountmodel.EntityName, err)
	}

	//If record is found in cache, return it
	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetAccount(ctx, cond, morekeys...)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
		}

		return nil, common.ErrCannotGetEntity(accountmodel.EntityName, err)
	}

	if record != nil {
		var paging paggingcommon.Paging
		paging.Process()

		createRecord := record.ToCreateAccount()

		key := cacheutil.GenerateKey(accountmodel.EntityName, cond, paging, commonfilter.Filter{})
		err := biz.rdbStore.SaveAccount(ctx, createRecord, key)

		if err != nil {
			return nil, common.ErrCannotCreateEntity(accountmodel.EntityName, err)
		}
	}

	return record, nil
}
