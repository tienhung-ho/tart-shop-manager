package authbusiness

import (
	"context"
	"errors"
	"os"
	"strconv"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/model/sql/account"
	hashutil "tart-shop-manager/internal/util/hash"
	tokenutil "tart-shop-manager/internal/util/token"
	"time"
)

type LoginStore interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
}

type loginBusiness struct {
	jwtService *jwtService
	store      LoginStore
}

func NewLoginServiceBiz(jwtService *jwtService, store LoginStore) *loginBusiness {
	return &loginBusiness{jwtService: jwtService, store: store}
}

func (biz *loginBusiness) Login(ctx context.Context, login *accountmodel.LoginAccount, morekeys ...string) (*accountmodel.Account, *tokenutil.Token, error) {

	record, err := biz.store.GetAccount(ctx, map[string]interface{}{"email": login.Email})

	if err != nil {
		return nil, nil, common.ErrEmailInvalid(accountmodel.EntityName, err)
	}

	cost := os.Getenv("COST")

	costInt, err := strconv.Atoi(cost)

	if err != nil {
		return nil, nil, common.ErrInternal(err)
	}

	newHash := hashutil.NewPasswordManager(costInt)

	if ok := newHash.VerifyPassword(record.Password, login.Password); !ok {
		return nil, nil, common.ErrPasswordInvalid(accountmodel.EntityName, errors.New("incorrect password"))
	}

	simpleRecord := record.ToSimpleAccount()

	timeExpireAccess := time.Duration(1 * time.Hour)
	accessToken, err := biz.jwtService.GenerateToken(record.AccountID, record.RoleID, record.Email, timeExpireAccess)

	if err != nil {
		return nil, nil, common.ErrInternal(err)
	}

	timeExpireRefresh := time.Duration(30 * 24 * time.Hour)
	refreshToken, err := biz.jwtService.GenerateToken(record.AccountID, record.RoleID, record.Email, timeExpireRefresh)

	token := tokenutil.NewToken(accessToken, refreshToken)

	return simpleRecord, token, nil
}
