package accountbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	casbinbusiness "tart-shop-manager/internal/service/policies"
	rolebusiness "tart-shop-manager/internal/service/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error)
}

type UpdateAccountCache interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	DeleteAccount(ctx context.Context, morekeys ...string) error
}

type updateAccountBusiness struct {
	store     UpdateAccountStorage
	roleStore rolebusiness.GetRoleStorage
	cache     UpdateAccountCache
	auth      casbinbusiness.Authorization
}

func NewUpdateAccount(store UpdateAccountStorage, roleStore rolebusiness.GetRoleStorage, cache UpdateAccountCache, authorization casbinbusiness.Authorization) *updateAccountBusiness {
	return &updateAccountBusiness{store, roleStore, cache, authorization}
}

func (biz *updateAccountBusiness) UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error) {

	record, err := biz.store.GetAccount(ctx, cond)

	if err != nil {
		return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateAccount(ctx, map[string]interface{}{"account_id": record.AccountID}, account, morekeys...)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}
	var pagging paggingcommon.Paging
	pagging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: accountmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(accountmodel.EntityName, err)
	}

	if err := biz.cache.DeleteAccount(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	if record.RoleID != updatedRecord.RoleID {
		// Xóa vai trò cũ của người dùng
		if err := biz.auth.RemoveUserFromAllRoles(ctx, record.Email); err != nil {
			return nil, common.ErrCannotDeleteEntity("user roles", err)
		}

		role, err := biz.roleStore.GetRole(ctx, map[string]interface{}{"role_id": updatedRecord.RoleID})

		if err != nil {
			return nil, common.ErrNotFoundEntity(rolemodel.EntityName, err)
		}

		if err := biz.auth.AddUserToRole(ctx, record.Email, role.Name); err != nil {
			return nil, common.ErrCannotCreateEntity("user roles", err)
		}

	}

	return updatedRecord, nil
}
