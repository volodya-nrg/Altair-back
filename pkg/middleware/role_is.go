package middleware

import (
	"altair/configs"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

// RoleIs - посредник, определяющий роль пользователя
func RoleIs() gin.HandlerFunc {
	return func(c *gin.Context) {
		authVal := c.Request.Header.Get("Authorization")
		bearerPrefix := "Bearer "
		pattern := `^` + bearerPrefix + `.+$`

		c.Set("roleIs", "")

		if matched, err := regexp.Match(pattern, []byte(authVal)); err != nil || !matched {
			c.Next()
			return
		}

		tokenSrc := authVal[len(bearerPrefix):]
		serviceSession := service.NewSessionService()
		accessTokenInfo, err := serviceSession.ParseAccessToken(tokenSrc)
		if err != nil || !accessTokenInfo.Verify(configs.Cfg.TokenPassword) {
			c.Next()
			return
		}

		diffTime := int(time.Until(time.Unix(accessTokenInfo.Exp, 0)) / time.Second)
		if diffTime < 1 {
			c.Next()
			return
		}

		c.Set("roleIs", accessTokenInfo.UserRole)
		c.Next()
	}
}
