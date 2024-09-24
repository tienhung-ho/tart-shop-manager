package authmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"os"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	rolestorage "tart-shop-manager/internal/repository/mysql/role"
	rolecache "tart-shop-manager/internal/repository/redis/role"
	authbusiness "tart-shop-manager/internal/service/auth"
	rolebusiness "tart-shop-manager/internal/service/role"
)

func AuthRequire(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
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

		store := rolestorage.NewMySQLRole(db)
		cache := rolecache.NewRdbStorage(rdb)
		biz := rolebusiness.NewGetBusinessBiz(store, cache)

		record, err := biz.GetRole(c, map[string]interface{}{"role_id": role})

		if err != nil {
			c.JSON(http.StatusNotFound, common.ErrEntityDeleted(rolemodel.EntityName, err))
			c.Abort()
			return
		}

		roleStatus := record.Status
		
		// Kiểm tra trạng thái vai trò
		switch *roleStatus {
		case common.StatusActive:
			// Trạng thái hợp lệ, tiếp tục
		case common.StatusInactive:
			c.JSON(http.StatusForbidden, common.NewUnauthorized(nil, "Your role is inactive. Please contact support.", "ErrRoleStatus", "ROLE_STATUS"))
			c.Abort()
			return
		case common.StatusPending:
			c.JSON(http.StatusForbidden, common.NewUnauthorized(nil, "Your role is pending approval.", "ErrRoleStatus", "ROLE_STATUS"))
			c.Abort()
			return
		default:
			c.JSON(http.StatusForbidden, common.NewUnauthorized(nil, "Your role status does not permit access.", "ErrRoleStatus", "ROLE_STATUS"))
			c.Abort()
			return
		}

		c.Set("id", id)
		c.Set("role", role)
		c.Set("email", email)

		c.Next()
	}
}
