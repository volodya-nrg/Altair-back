package middleware

import (
	"altair/configs"
	"altair/pkg/manager"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"regexp"
	"time"
)

// Auth - посредник, в котором проверяется (не)зарегистрированный пользователь
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Next()
		authVal := c.Request.Header.Get("Authorization")
		bearerPrefix := "Bearer "
		pattern := `^` + bearerPrefix + `.+$`

		if matched, err := regexp.Match(pattern, []byte(authVal)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			}) // 401
			return

		} else if !matched {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": manager.ErrNotMatched.Error(),
			}) // 401
			return
		}

		tokenSrc := authVal[len(bearerPrefix):]
		serviceSession := service.NewSessionService()

		accessTokenInfo, err := serviceSession.ParseAccessToken(tokenSrc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			}) // 401
			return
		}

		if !accessTokenInfo.Verify(configs.Cfg.TokenPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": manager.ErrNotVerify.Error(),
			}) // 401
			return
		}

		// время жизни access-token-а
		diffTime := int(time.Until(time.Unix(accessTokenInfo.Exp, 0)) / time.Second)

		// если время жизни access-token-а прошло, то надо посмотреть - есть ли кука и актуальна ли сессия.
		// Если все норм, то пропускаем через MiddleWare
		if diffTime < 1 {
			cookieValue, err := c.Cookie(manager.CookieTokenName)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusPreconditionRequired, err.Error()) // 428 Необходимо предусловие
				return

			} else if cookieValue == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, manager.ErrNotFoundCookie.Error()) // 401
				return
			}

			session, err := serviceSession.GetSessionByRefreshToken(cookieValue)
			if gorm.IsRecordNotFoundError(err) || err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error()) // 401
				return
			}

			if err := serviceSession.Delete(session.SessionID, nil); err != nil {
				c.AbortWithStatusJSON(500, err.Error())
				return
			}

			// если сессия еще действует, то идем далее по коду (чтоб пересоздать ее). Иначе просим авторизоваться.
			timeDiff := int(time.Until(session.ExpiresIn) / time.Second)
			if timeDiff < 1 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": manager.ErrSessionIsOver.Error(),
				}) // 401
				return
			}
		}

		c.Set("userID", accessTokenInfo.UserID)
		c.Set("userRole", accessTokenInfo.UserRole)

		c.Next()
	}
}
