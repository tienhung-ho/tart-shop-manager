package authmodel

import "github.com/golang-jwt/jwt/v5"

type AccountJwtClaims struct {
	AccountId uint64 `json:"account_id"`
	Role      uint8  `json:"role"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}
