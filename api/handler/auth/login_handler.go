package authhandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	accountstorage "tart-shop-manager/internal/repository/mysql/account"
	authbusiness "tart-shop-manager/internal/service/auth"
)

func LoginHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var login accountmodel.LoginAccount

		if err := c.ShouldBind(&login); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrCanNotBindEntity(accountmodel.EntityName, err))
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		store := accountstorage.NewMySQLAccount(db)
		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)
		biz := authbusiness.NewLoginServiceBiz(jwtService, store)

		simpleUser, token, err := biz.Login(c.Request.Context(), &login)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.SetCookie("access_token", token.AccessToken, 3600, "/", "", true, true)

		// Lưu Refresh Token vào Cookie
		c.SetCookie("refresh_token", token.RefreshToken, 30*24*3600, "/", "", true, true)

		c.JSON(http.StatusOK, common.NewReponseUserToken(token.AccessToken, token.RefreshToken, simpleUser))

	}
}
