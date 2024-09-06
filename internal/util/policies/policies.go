package casbinutil

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"sync"
)

var (
	instance *casbin.Enforcer
	once     sync.Once
)

// InitEnforcer initializes the Casbin enforcer as a singleton.
func InitEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		var adapter *gormadapter.Adapter
		adapter, err = gormadapter.NewAdapterByDB(db)
		if err != nil {
			return
		}
		instance, err = casbin.NewEnforcer(modelPath, adapter)
		if err != nil {
			return
		}
		err = instance.LoadPolicy()
	})
	return instance, err
}
