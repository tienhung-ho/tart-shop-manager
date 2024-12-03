package ingredientbusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"

	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemSotrage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]ingredientmodel.Ingredient, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, key string) ([]ingredientmodel.Ingredient, error)
	SaveIngredient(ctx context.Context, data interface{}, morekeys ...string) error
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
	SaveFilter(ctx context.Context, filter *commonfilter.Filter, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
}

type listItemIngredientBusiness struct {
	store ListItemSotrage
	cache ListItemCache
}

func NewListItemIngredientBiz(store ListItemSotrage, cache ListItemCache) *listItemIngredientBusiness {
	return &listItemIngredientBusiness{store, cache}
}

func (biz *listItemIngredientBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]ingredientmodel.Ingredient, error) {

	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	baseKey, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: ingredientmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", ingredientmodel.EntityName),
	})

	ingredientKey := baseKey
	pagingKey := baseKey + ":paging"
	if err != nil {
		return nil, common.ErrCannotGenerateKey(ingredientmodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, ingredientKey)

	if err != nil {
		return nil, common.ErrCannotListEntity(ingredientmodel.EntityName, err)
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

			return nil, common.ErrNotFoundEntity(ingredientmodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(ingredientmodel.EntityName, err)
	}

	if len(records) != 0 {

		if err := biz.cache.SaveIngredient(ctx, records, ingredientKey); err != nil {
			return nil, common.ErrCannotCreateEntity(ingredientmodel.EntityName, err)
		}

		if err := biz.cache.SavePaging(ctx, paging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(ingredientmodel.EntityName, err)
		}
	}

	return records, nil
}
