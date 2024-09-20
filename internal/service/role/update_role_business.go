package rolebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	casbinbusiness "tart-shop-manager/internal/service/policies"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateRoleStorage interface {
	GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error)
	UpdateRole(ctx context.Context, cond map[string]interface{}, data *rolemodel.UpdateRole, morekeys ...string) error
}

type UpdateRoleCache interface {
	GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error)
	DeleteRole(ctx context.Context, morekeys ...string) error
}

type updateRoleBusiness struct {
	store    UpdateRoleStorage
	cache    UpdateRoleCache
	perStore ListPermissionStorage
	auth     casbinbusiness.Authorization
}

func NewUpdateRolebiz(store UpdateRoleStorage, cache UpdateRoleCache, perStore ListPermissionStorage, auth casbinbusiness.Authorization) *updateRoleBusiness {
	return &updateRoleBusiness{store, cache, perStore, auth}
}

func (biz *updateRoleBusiness) UpdateRole(ctx context.Context, cond map[string]interface{}, data *rolemodel.UpdateRole, morekeys ...string) error {

	record, err := biz.store.GetRole(ctx, cond, morekeys...)

	if err != nil {
		return common.ErrNotFoundEntity(rolemodel.EntityName, err)
	}

	// Lấy tất cả các id permission từ dữ liệu đầu vào
	var permissionIds []uint
	for _, perm := range data.Permissions {
		if perm.PermissionID != 0 {
			permissionIds = append(permissionIds, perm.PermissionID)
		}
	}

	// Tìm tất cả các permissions tồn tại trong database
	var condi map[string]interface{}

	if len(permissionIds) > 0 {
		condi = map[string]interface{}{
			"permission_id": permissionIds,
		}
	}

	var paging paggingcommon.Paging

	paging.Process()

	permissions, err := biz.perStore.ListItem(ctx, condi, &paging, &commonfilter.Filter{})
	if err != nil {
		return err
	}

	role := rolemodel.UpdateRole{
		Name:         data.Name,
		Description:  data.Description,
		Permissions:  permissions,
		CommonFields: data.CommonFields,
	}

	if err := biz.store.UpdateRole(ctx, map[string]interface{}{"role_id": record.RoleID}, &role, morekeys...); err != nil {
		return common.ErrCannotUpdateEntity(rolemodel.EntityName, err)
	}

	if err := biz.auth.RemoveAllPolicesOfRole(record.Name); err != nil {
		return err
	}

	// Add policies for the role using Authorization interface
	if err := biz.auth.AddPoliciesForRole(role.Name, role.Permissions); err != nil {
		return err
	}

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: rolemodel.EntityName,
		Cond:       cond,
		Paging:     paging,
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
