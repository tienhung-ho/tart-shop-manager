package casbinbusiness

import (
	"context"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
)

type Authorization interface {
	AddPoliciesForRole(roleName string, permissions []permissionmodel.Permission) error
	RemoveUserFromAllRoles(ctx context.Context, user string) error
	AddUserToRole(ctx context.Context, user string, role string) error
}
