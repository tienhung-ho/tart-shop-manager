package accountbusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]accountmodel.Account, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, key string) ([]accountmodel.Account, error)
	SaveAccount(ctx context.Context, data interface{}, morekeys ...string) error
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
}

type listItemBusiness struct {
	store ListItemStorage
	cache ListItemCache
}

func NewListItemBiz(store ListItemStorage, cache ListItemCache) *listItemBusiness {
	return &listItemBusiness{store, cache}
}

func (biz *listItemBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]accountmodel.Account, error) {

	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *paging
	filterCopy := *filter

	baseKey, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: accountmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", accountmodel.EntityName),
	})

	accountKey := baseKey
	pagingKey := baseKey + ":paging"

	if err != nil {
		return nil, common.ErrCannotGenerateKey(accountmodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, accountKey)

	if err != nil {
		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}

	if len(records) != 0 {
		cachedPaging, err := biz.cache.GetPaging(ctx, pagingKey)
		if err == nil {
			paging.Page = cachedPaging.Page
			paging.Total = cachedPaging.Total
			paging.Limit = cachedPaging.Limit
			paging.Sort = cachedPaging.Sort
		}
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
		if err := biz.cache.SaveAccount(ctx, records, accountKey); err != nil {
			return nil, common.ErrCannotCreateEntity(accountmodel.EntityName, err)
		}
		if err := biz.cache.SavePaging(ctx, paging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
	}

	return records, nil
}
