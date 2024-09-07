package accountbusiness

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	casbinbusiness "tart-shop-manager/internal/service/policies"
	rolebusiness "tart-shop-manager/internal/service/role"
	hashutil "tart-shop-manager/internal/util/hash"
	responseutil "tart-shop-manager/internal/util/response"
)

type CreateAccountBusiness interface {
	CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error)
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
}

type createAccountBusiness struct {
	store     CreateAccountBusiness
	roleStore rolebusiness.GetRoleStorage
	auth      casbinbusiness.Authorization
}

func NewCreateAccountbiz(store CreateAccountBusiness, roleStore rolebusiness.GetRoleStorage, auth casbinbusiness.Authorization) *createAccountBusiness {
	return &createAccountBusiness{store: store, roleStore: roleStore, auth: auth}
}

func (biz *createAccountBusiness) CreateAccount(ctx context.Context, data *accountmodel.CreateAccount, morekeys ...string) (uint64, error) {

	costEnv := os.Getenv("COST")
	costInt, err := strconv.Atoi(costEnv)

	if err != nil {
		return 0, err
	}

	hashUtil := hashutil.NewPasswordManager(costInt)
	data.Password, err = hashUtil.HashPassword(data.Password)

	if err != nil {
		return 0, err
	}

	recordId, err := biz.store.CreateAccount(ctx, data, morekeys...)

	if err != nil {
		// Check for MySQL duplicate entry error

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, accountmodel.EntityName) // Extract field causing the duplicate error
			return 0, common.ErrDuplicateEntry(accountmodel.EntityName, fieldName, err)
		}

		return 0, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	if data.RoleID != 0 {
		role, err := biz.roleStore.GetRole(ctx, map[string]interface{}{"role_id": data.RoleID})

		if err != nil {
			return 0, common.ErrNotFoundEntity(rolemodel.EntityName, err)
		}

		if err := biz.auth.AddUserToRole(ctx, data.Email, role.Name); err != nil {
			return 0, common.ErrCannotCreateEntity("user roles", err)
		}

	}

	return recordId, nil
}
