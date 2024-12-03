package permissionbusiness

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemPermissionStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]permissionmodel.Permission, error)
}

type ListItemPermissionCache interface {
	SavePermission(ctx context.Context, data interface{}, morekeys ...string) error
	ListItem(ctx context.Context, key string) ([]permissionmodel.Permission, error)
}

type listItemPermissionBusiness struct {
	store ListItemPermissionStorage
	cache ListItemPermissionCache
}

func NewListItemPermissionBiz(store ListItemPermissionStorage, cache ListItemPermissionCache) *listItemPermissionBusiness {
	return &listItemPermissionBusiness{store, cache}
}

func (biz *listItemPermissionBusiness) ListItemPermission(ctx context.Context, cond map[string]interface{},
	paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]permissionmodel.Permission, error) {

	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *paging
	filterCopy := *filter

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: permissionmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
		KeyType:    fmt.Sprintf("List:%s:", permissionmodel.EntityName),
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(permissionmodel.EntityName, err)
	}

	records, err := biz.cache.ListItem(ctx, key)

	if err != nil {
		log.Print(err)
		return nil, common.ErrCannotListEntity(permissionmodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, paging, filter, morekeys...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(permissionmodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(permissionmodel.EntityName, err)
	}

	if len(records) != 0 {

		if err := biz.cache.SavePermission(ctx, records, key); err != nil {
			return nil, common.ErrCannotUpdateEntity(permissionmodel.EntityName, err)
		}

	}

	return records, nil
}
