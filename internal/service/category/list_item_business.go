package categorybusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemSotrage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]categorymodel.Category, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, key string) ([]categorymodel.Category, error)
	SaveCategory(ctx context.Context, data interface{}, morekeys ...string) error
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
	SaveFilter(ctx context.Context, filter *commonfilter.Filter, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
}

type listItemCategoryBusiness struct {
	store ListItemSotrage
	cache ListItemCache
}

func NewListItemCategoryBiz(store ListItemSotrage, cache ListItemCache) *listItemCategoryBusiness {
	return &listItemCategoryBusiness{store, cache}
}

func (biz *listItemCategoryBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]categorymodel.Category, error) {

	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: categorymodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", categorymodel.EntityName),
	})

	categoryKey := key
	pagingKey := key + ":paging"
	if err != nil {
		return nil, common.ErrCannotGenerateKey(categorymodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, categoryKey)

	if err != nil {
		return nil, common.ErrCannotListEntity(categorymodel.EntityName, err)
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

			return nil, common.ErrNotFoundEntity(categorymodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(categorymodel.EntityName, err)
	}

	if len(records) != 0 {

		if err := biz.cache.SaveCategory(ctx, records, key); err != nil {
			return nil, common.ErrCannotCreateEntity(categorymodel.EntityName, err)
		}

		if err := biz.cache.SavePaging(ctx, paging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(categorymodel.EntityName, err)
		}
	}

	return records, nil
}
