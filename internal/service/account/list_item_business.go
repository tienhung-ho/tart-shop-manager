package accountbusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]accountmodel.Account, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]accountmodel.Account, error)
	SaveAccount(ctx context.Context, data interface{}, morekeys ...string) error
}

type listItemBusiness struct {
	store ListItemStorage
	cache ListItemCache
}

func NewListItemBiz(store ListItemStorage, cache ListItemCache) *listItemBusiness {
	return &listItemBusiness{store, cache}
}

func (biz *listItemBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]accountmodel.Account, error) {

	records, err := biz.cache.ListItem(ctx, cond, paging, filter, morekeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, paging, filter, morekeys...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}

	if len(records) != 0 {

		key := cacheutil.GenerateKey(accountmodel.EntityName, cond, *paging, *filter)
		err := biz.cache.SaveAccount(ctx, records, key)

		if err != nil {
			return nil, common.ErrCannotCreateEntity(accountmodel.EntityName, err)
		}
	}

	return records, nil
}
