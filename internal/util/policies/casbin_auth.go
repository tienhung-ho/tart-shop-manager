package casbinutil

import (
	"github.com/casbin/casbin/v2"
	"log"
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
