package authmiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"tart-shop-manager/internal/common"
	authbusiness "tart-shop-manager/internal/service/auth"
)

func AuthRequire() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, common.NewUnauthorized(err, "This action requires login to perform", "ErrRequireLogin", "ACCESS_TOKEN"))
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)

		claims, err := jwtService.ValidateToken(accessToken)

		if err != nil {
			switch err.Error() {
			case "token is expired":
				c.JSON(http.StatusUnauthorized, common.TokenExpired("Access Token", err))
			case "token signature is invalid":
				c.JSON(http.StatusUnauthorized, common.NewUnauthorized(err, "Token signature is invalid", "ErrInvalidTokenSignature", "ACCESS_TOKEN"))
			case "invalid token issuer":
				c.JSON(http.StatusUnauthorized, common.NewUnauthorized(err, "Token issuer is invalid", "ErrInvalidTokenIssuer", "ACCESS_TOKEN"))
			case "invalid token audience":
				c.JSON(http.StatusUnauthorized, common.NewUnauthorized(err, "Token audience is invalid", "ErrInvalidTokenAudience", "ACCESS_TOKEN"))
			default:
				c.JSON(http.StatusUnauthorized, common.ErrInternal(err))
			}
			c.Abort()
			return
		}

		id := claims.AccountId
		role := claims.Role
		email := claims.Email

		c.Set("id", id)
		c.Set("role", role)
		c.Set("email", email)

		c.Next()
	}
}
