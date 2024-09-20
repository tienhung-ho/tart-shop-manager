package rolebusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemRoleStorage interface {
	ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error)
}

type ListItemRoleCache interface {
	SaveRole(ctx context.Context, data interface{}, morekeys ...string) error
	ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error)
}

type listItemRoleBusiness struct {
	store ListItemRoleStorage
	cache ListItemRoleCache
}

func NewListItemRoleBiz(store ListItemRoleStorage, cache ListItemRoleCache) *listItemRoleBusiness {
	return &listItemRoleBusiness{store, cache}
}

func (biz *listItemRoleBusiness) ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error) {

	records, err := biz.cache.ListItemRole(ctx, cond, paging, filter, morekeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(rolemodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItemRole(ctx, cond, paging, filter, morekeys...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}

	if len(records) != 0 {

		// Generate cache key
		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: rolemodel.EntityName,
			Cond:       cond,
			Paging:     *paging,
			Filter:     *filter,
			MoreKeys:   morekeys,
		})
		if err != nil {
			return nil, common.ErrCannotGenerateKey(rolemodel.EntityName, err)
		}

		if err := biz.cache.SaveRole(ctx, records, key); err != nil {
			return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
		}

	}

	return records, nil
}
