package authmiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	casbinutil "tart-shop-manager/internal/util/policies"
)

// CasbinMiddleware kiểm tra quyền truy cập dựa trên Casbin policies
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		enforcer := casbinutil.GetEnforcer()

		if enforcer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Casbin enforcer not initialized"})
			c.Abort()
			return
		}
		userEmail := c.GetString("email") // Lấy vai trò của người dùng từ context
		path := c.Request.URL.Path        // Đường dẫn hiện tại
		method := c.Request.Method        // Phương thức HTTP (GET, POST, DELETE, ...)

		// Kiểm tra quyền của người dùng với Casbin

		ok, err := enforcer.Enforce(userEmail, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next() // Quyền hợp lệ, tiếp tục xử lý yêu cầu
	}
}
