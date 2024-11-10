package authbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	authmodel "tart-shop-manager/internal/entity/dtos/auth"
	tokenutil "tart-shop-manager/internal/util/token"
	"time"
)

type refreshTokenBusiness struct {
	jwtService *jwtService
}

func NewRefreshTokenBusiness(jwtService *jwtService) *refreshTokenBusiness {
	return &refreshTokenBusiness{
		jwtService: jwtService,
	}
}

func (biz *refreshTokenBusiness) RefreshToken(ctx context.Context, claim *authmodel.AccountJwtClaims) (*tokenutil.Token, error) {

	timeExpireAccess := time.Duration(1 * time.Hour)
	accessToken, err := biz.jwtService.GenerateToken(claim.AccountId, claim.Role, claim.Email, timeExpireAccess)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	timeExpireRefresh := time.Duration(30 * 1 * time.Hour)
	refreshToken, err := biz.jwtService.GenerateToken(claim.AccountId, claim.Role, claim.Email, timeExpireRefresh)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return tokenutil.NewToken(accessToken, refreshToken), nil
}
