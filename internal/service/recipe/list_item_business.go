package recipebusiness

import (
	"context"
	"fmt"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemRecipeStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]recipemodel.Recipe, error)
}

type ListItemRecipeCache interface {
	ListItem(ctx context.Context, key string) ([]recipemodel.Recipe, error)
	SaveRecipe(ctx context.Context, data interface{}, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
}

type newListItemRecipeBusiness struct {
	store ListItemRecipeStorage
	cache ListItemRecipeCache
}

func NewListItemRecipeBiz(store ListItemRecipeStorage, cache ListItemRecipeCache) *newListItemRecipeBusiness {
	return &newListItemRecipeBusiness{store, cache}
}

func (biz *newListItemRecipeBusiness) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]recipemodel.Recipe, error) {

	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *paging
	filterCopy := *filter

	// Tạo khóa cache
	baseKey, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: recipemodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", recipemodel.EntityName),
	})

	productKey := baseKey
	pagingKey := baseKey + ":paging"

	if err != nil {
		return nil, common.ErrCannotGenerateKey(recipemodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, productKey)

	if err != nil {
		return nil, common.ErrCannotListEntity(recipemodel.EntityName, err)
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
		return nil, common.ErrCannotListEntity(recipemodel.EntityName, err)
	}

	if len(records) != 0 {
		if err := biz.cache.SaveRecipe(ctx, records, productKey); err != nil {
			return nil, common.ErrCannotCreateEntity(recipemodel.EntityName, err)
		}
		if err := biz.cache.SavePaging(ctx, paging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
	}

	return records, nil
}
