package rolebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteRoleStorage interface {
	GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error)
	DeleteRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteRoleCache interface {
	DeleteRole(ctx context.Context, morekeys ...string) error
}

type deleteRoleBusiness struct {
	store DeleteRoleStorage
	cache DeleteRoleCache
}

func NewDeleteRoleBiz(store DeleteRoleStorage, cache DeleteRoleCache) *deleteRoleBusiness {
	return &deleteRoleBusiness{store, cache}
}

func (biz *deleteRoleBusiness) DeleteRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetRole(ctx, cond, morekeys...)

	if err != nil {
		return common.ErrNotFoundEntity(rolemodel.EntityName, err)
	}

	if err := biz.store.DeleteRole(ctx, map[string]interface{}{"role_id": record.RoleID}, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(rolemodel.EntityName, err)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: rolemodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return common.ErrCannotGenerateKey(rolemodel.EntityName, err)
	}

	if err := biz.cache.DeleteRole(ctx, key); err != nil {
		return common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	return nil
}
