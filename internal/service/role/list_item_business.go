package rolebusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemRoleStorage interface {
	ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error)
}

type ListItemRoleCache interface {
	SaveRole(ctx context.Context, data interface{}, morekeys ...string) error
	ListItem(ctx context.Context, key string) ([]rolemodel.Role, error)
	SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error)
}

type listItemRoleBusiness struct {
	store ListItemRoleStorage
	cache ListItemRoleCache
}

func NewListItemRoleBiz(store ListItemRoleStorage, cache ListItemRoleCache) *listItemRoleBusiness {
	return &listItemRoleBusiness{store, cache}
}

func (biz *listItemRoleBusiness) ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error) {

	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *paging
	filterCopy := *filter

	// Generate cache key
	baseKey, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: rolemodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", rolemodel.EntityName),
	})

	roleKey := baseKey
	pagingKey := baseKey + ":paging"
	if err != nil {
		return nil, common.ErrCannotGenerateKey(rolemodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, roleKey)

	if err != nil {
		log.Print(err)
		return nil, common.ErrCannotListEntity(rolemodel.EntityName, err)
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

	records, err = biz.store.ListItemRole(ctx, cond, paging, filter, morekeys...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(rolemodel.EntityName, err)
		}
		//log.Print(err)
		return nil, common.ErrCannotListEntity(rolemodel.EntityName, err)
	}

	if len(records) != 0 {

		if err := biz.cache.SaveRole(ctx, records, roleKey); err != nil {
			return nil, common.ErrCannotUpdateEntity(rolemodel.EntityName, err)
		}

		if err := biz.cache.SavePaging(ctx, paging, pagingKey); err != nil {
			return nil, common.ErrCannotCreateEntity(rolemodel.EntityName, err)
		}

	}

	return records, nil
}
