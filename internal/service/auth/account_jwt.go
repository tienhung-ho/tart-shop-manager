package authbusiness

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	authmodel "tart-shop-manager/internal/entity/dtos/auth"
	"time"
)

type JwtService interface {
	GenerateToken(id uint64, roleId int, email string, expireTime int) (string, error)
	ValidateToken(tokenString string) (*authmodel.AccountJwtClaims, error)
}

type jwtService struct {
	secretkey string
	issuer    string
	audience  string
}

func NewJwtService(secretkey, issuer, audience string) *jwtService {
	return &jwtService{secretkey, issuer, audience}
}

func (s *jwtService) GenerateToken(id uint64, roleId uint8, email string, expireTime time.Duration) (string, error) {

	claims := &authmodel.AccountJwtClaims{
		AccountId: id,
		Role:      roleId,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{s.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	signedToken, err := token.SignedString([]byte(s.secretkey))
	if err != nil {
		return "", errors.New("error signing token")
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*authmodel.AccountJwtClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &authmodel.AccountJwtClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secretkey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("token signature is invalid")
		}

		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token is not valid yet")
		}

		return nil, errors.New("error parsing token")

	}

	claims, ok := token.Claims.(*authmodel.AccountJwtClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Kiểm tra thêm Issuer và Audience
	if claims.Issuer != s.issuer {
		return nil, errors.New("invalid token issuer")
	}

	// Kiểm tra Audience thủ công
	if len(claims.Audience) == 0 || claims.Audience[0] != s.audience {
		return nil, errors.New("invalid token audience")
	}

	return claims, nil
}
