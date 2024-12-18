package rolebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	casbinbusiness "tart-shop-manager/internal/service/policies"
)

type CreateRoleStorage interface {
	CreateRole(ctx context.Context, data *rolemodel.CreateRole, morekeys ...string) (uint, error)
}

type CreateRoleCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type ListPermissionStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]permissionmodel.Permission, error)
}

type createRoleBusiness struct {
	store    CreateRoleStorage
	perStore ListPermissionStorage
	cache    CreateRoleCache
	auth     casbinbusiness.Authorization
}

func NewCreateRoleBiz(store CreateRoleStorage, perStore ListPermissionStorage, cache CreateRoleCache, auth casbinbusiness.Authorization) *createRoleBusiness {
	return &createRoleBusiness{store: store, perStore: perStore, cache: cache, auth: auth}
}

func (biz *createRoleBusiness) CreateRole(ctx context.Context, data *rolemodel.CreateRole, morekeys ...string) (uint, error) {

	// Lấy tất cả các tên permission từ dữ liệu đầu vào
	//var permissionNames []string
	var permissionIds []uint
	for _, perm := range data.Permissions {
		// Kiểm tra nếu Name không rỗng và PermissionID khác 0 thì mới thêm vào slice
		//if perm.Name != "" {
		//	permissionNames = append(permissionNames, perm.Name)
		//}
		if perm.PermissionID != 0 {
			permissionIds = append(permissionIds, perm.PermissionID)
		}
	}

	// Tìm tất cả các permissions tồn tại trong database
	var cond map[string]interface{}

	if len(permissionIds) > 0 {
		cond = map[string]interface{}{
			"permission_id": permissionIds,
		}
	}

	var paging paggingcommon.Paging

	paging.Process()

	permissions, err := biz.perStore.ListItem(ctx, cond, &paging, &commonfilter.Filter{})
	if err != nil {
		return 0, err
	}

	role := rolemodel.CreateRole{
		Name:         data.Name,
		Description:  data.Description,
		Permissions:  permissions,
		CommonFields: data.CommonFields,
	}

	recordId, err := biz.store.CreateRole(ctx, &role)

	if err != nil {
		return 0, common.ErrCannotCreateEntity(rolemodel.EntityName, err)
	}

	// Add policies for the role using Authorization interface
	err = biz.auth.AddPoliciesForRole(role.Name, role.Permissions)
	if err != nil {
		return 0, err
	}

	if err := biz.cache.DeleteListCache(ctx, rolemodel.EntityName); err != nil {
		return 0, common.ErrCannotCreateEntity(rolemodel.EntityName, err)
	}

	return recordId, nil
}
