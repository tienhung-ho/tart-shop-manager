package accountbusiness

import (
	"context"
	"errors"
	"tart-shop-manager/internal/common"
	accountmodel2 "tart-shop-manager/internal/entity/model/sql/account"
)

type CreateAccountBusiness interface {
	CreateAccount(ctx context.Context, data *accountmodel2.CreateAccount, morekeys ...string) (uint64, error)
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel2.Account, error)
}

type createAccountBusiness struct {
	store CreateAccountBusiness
}

func NewCreateAccountbiz(store CreateAccountBusiness) *createAccountBusiness {
	return &createAccountBusiness{store: store}
}

func (biz *createAccountBusiness) CreateAccount(ctx context.Context, data *accountmodel2.CreateAccount, morekeys ...string) (uint64, error) {

	// Check if the email already exists
	existingByEmail, err := biz.store.GetAccount(ctx, map[string]interface{}{"email": data.Email}, morekeys...)
	if err == nil && existingByEmail != nil {
		return 0, common.ErrDuplicateEntry(accountmodel2.EntityName, "email", errors.New("email already exists"))
	}

	// Check if the phone already exists
	existingByPhone, err := biz.store.GetAccount(ctx, map[string]interface{}{"phone": data.Phone}, morekeys...)
	if err == nil && existingByPhone != nil {
		return 0, common.ErrDuplicateEntry(accountmodel2.EntityName, "phone", errors.New("phone already exists"))
	}

	recordId, err := biz.store.CreateAccount(ctx, data, morekeys...)

	if err != nil {
		return 0, common.ErrCannotCreateEntity(accountmodel2.EntityName, err)
	}

	return recordId, nil
}
