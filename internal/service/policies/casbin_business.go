package casbinbusiness

import permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"

type Authorization interface {
	AddPoliciesForRole(roleName string, permissions []permissionmodel.Permission) error
}
