package middleware

import (
	"altair/pkg/manager"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Admin - посредник в котором осуществляется проверка на (не)админ
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.MustGet("userRole").(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, manager.ErrAccessDined.Error()) // 403
			return
		}

		if userRole != manager.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, manager.ErrAccessDined.Error()) // 403
			return
		}

		c.Next()
	}
}
