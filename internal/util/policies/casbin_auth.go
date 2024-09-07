package casbinutil

import (
	"context"
	"github.com/casbin/casbin/v2"
	"log"
	"tart-shop-manager/internal/common"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

// CasbinAuthorization struct để implement Authorization interface
type casbinAuthorization struct {
	enforcer *casbin.Enforcer
}

// NewCasbinAuthorization returns a new CasbinAuthorization instance
func NewCasbinAuthorization(enforcer *casbin.Enforcer) *casbinAuthorization {
	return &casbinAuthorization{enforcer: enforcer}
}

// AddPoliciesForRole thêm các policies dựa trên role và permissions của role
func (auth *casbinAuthorization) AddPoliciesForRole(roleName string, permissions []permissionmodel.Permission) error {
	for _, permission := range permissions {
		// Casbin policy format: (roleName, object, action)
		success, err := auth.enforcer.AddPolicy(roleName, permission.Object, permission.Action)
		if err != nil {
			return err
		}
		if !success {
			log.Printf("Policy already exists for role: %s, object: %s, action: %s", roleName, permission.Object, permission.Action)
		}
	}
	return nil
}

// AddUserToRole thêm một người dùng vào một nhóm quyền
func (biz *casbinAuthorization) AddUserToRole(ctx context.Context, user string, role string) error {
	// Kiểm tra nếu người dùng đã có vai trò này chưa
	hasRole, err := biz.enforcer.HasRoleForUser(user, role)
	if err != nil {
		return common.ErrCannotCreateEntity("user to role", err)
	}

	// Nếu người dùng chưa có vai trò này, thì thêm
	if !hasRole {
		_, err := biz.enforcer.AddRoleForUser(user, role)
		if err != nil {
			return common.ErrCannotCreateEntity("user to role", err)
		}
	}
	return nil
}

// RemoveUserFromRole xóa một người dùng khỏi một nhóm quyền
func (biz *casbinAuthorization) RemoveUserFromRole(ctx context.Context, user string, role string) error {
	// Xóa người dùng khỏi vai trò
	_, err := biz.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return common.ErrCannotDeleteEntity("user from role", err)
	}
	return nil
}

// GetRolesForUser lấy danh sách các vai trò của một người dùng
func (biz *casbinAuthorization) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	roles, err := biz.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, common.ErrCannotGetEntity("roles for user", err)
	}
	return roles, nil
}

// RemoveUserFromAllRoles xóa người dùng khỏi tất cả các vai trò
func (biz *casbinAuthorization) RemoveUserFromAllRoles(ctx context.Context, user string) error {
	// Xóa người dùng khỏi tất cả các vai trò
	_, err := biz.enforcer.DeleteRolesForUser(user)
	if err != nil {
		return common.ErrCannotDeleteEntity("user from all roles", err)
	}
	return nil
}
