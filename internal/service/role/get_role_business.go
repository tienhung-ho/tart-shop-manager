package rolebusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetRoleStorage interface {
	GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error)
}

type GetRoleCache interface {
	GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error)
	SaveRole(ctx context.Context, data interface{}, morekeys ...string) error
}

type getRoleBusiness struct {
	store GetRoleStorage
	cache GetRoleCache
}

func NewGetBusinessBiz(store GetRoleStorage, cache GetRoleCache) *getRoleBusiness {
	return &getRoleBusiness{store, cache}
}

func (biz *getRoleBusiness) GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error) {

	record, err := biz.cache.GetRole(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrCannotGetEntity(rolemodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetRole(ctx, cond, morekeys...)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
		}

		return nil, common.ErrCannotGetEntity(rolemodel.EntityName, err)
	}

	if record != nil {
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
			return nil, common.ErrCannotGenerateKey(rolemodel.EntityName, err)
		}

		createRole := record.ToCreateRoleCache()

		if err := biz.cache.SaveRole(ctx, createRole, key); err != nil {
			log.Print(err)
			return nil, common.ErrCannotCreateEntity(accountmodel.EntityName, err)
		}
	}

	return record, nil
}
