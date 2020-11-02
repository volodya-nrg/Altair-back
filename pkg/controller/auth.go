package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/configs"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/pkg/soc"
	"altair/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// GetAuthLogout - выход из профиля
func GetAuthLogout(c *gin.Context) {
	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserID.Error())
		return
	}
	//----------------------------------------------------------

	serviceSession := service.NewSessionService()

	if err := serviceSession.DeleteAllByUserID(userID, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	secure := false

	if configs.Cfg.Mode == gin.ReleaseMode {
		secure = true
	}

	c.SetCookie(manager.CookieTokenName, "", 0,
		manager.CookiePath, manager.CookieDomain, secure, true)

	c.JSON(http.StatusNoContent, nil) // 204
}

// PostAuthLogin - авторизация
func PostAuthLogin(c *gin.Context) {
	postRequest := new(request.PostLogin)
	serviceUsers := service.NewUserService()
	serviceSession := service.NewSessionService()
	servicePhone := service.NewPhoneService()
	isAuthorizedThroughSoc := false

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error()) // 400
		return
	}

	email := postRequest.Email
	password := postRequest.Password
	code := postRequest.Code
	state := postRequest.State

	// если пытаются зайти через е-мэйл соц. сети, то выдадим ошибку
	if manager.IsSocialEmail(email) {
		c.JSON(400, manager.ErrEmailNotCorrect.Error())
		return
	}

	// если авторазация происходит через соц. сети
	if code != "" && state != "" {
		if has, _ := manager.InArray(state, manager.AvailableKindSoc); !has {
			logger.Warning.Println(manager.ErrUndefinedOptSoc.Error())
			c.JSON(500, manager.ErrUndefinedOptSoc.Error())
			return
		}

		var abstractSoc soc.Socer

		switch state {
		case "vk":
			abstractSoc = soc.NewVk(configs.Cfg.Socials.Vk.ClientID, configs.Cfg.Socials.Vk.ClientSecret, code)
		case "ok":
			abstractSoc = soc.NewOk(configs.Cfg.Socials.Ok.ClientID, configs.Cfg.Socials.Ok.ClientSecret, configs.Cfg.Socials.Ok.ClientPublic, code)
		case "fb":
			abstractSoc = soc.NewFb(configs.Cfg.Socials.Fb.ClientID, configs.Cfg.Socials.Fb.ClientSecret, code)
		case "ggl":
			abstractSoc = soc.NewGgl(configs.Cfg.Socials.Ggl.ClientID, configs.Cfg.Socials.Ggl.ClientSecret, code)
		}

		commonUserInfo, err := soc.CommonHandler(state, abstractSoc)
		if err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}

		if commonUserInfo.Email == "" {
			logger.Warning.Println(manager.ErrUndefinedSocEmail.Error())
			c.JSON(500, manager.ErrUndefinedSocEmail.Error())
			return
		}

		_, err = serviceUsers.GetUserByEmail(commonUserInfo.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return

		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			tmpUser := new(storage.User)
			tmpUser.Email = commonUserInfo.Email
			tmpUser.Password = manager.HashAndSalt(manager.RandStringRunes(10))
			tmpUser.IsEmailConfirmed = true

			if err := serviceUsers.Create(tmpUser, nil); err != nil {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return
			}
		}

		email = commonUserInfo.Email
		isAuthorizedThroughSoc = true
	}

	// стандартный подход проверки пользователя
	user, err := serviceUsers.GetUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, manager.ErrNotCorrectLoginPassword.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	phones, err := servicePhone.GetPhonesByUserID(user.UserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// те кто заходит через соц. сети не проверяем у них пароль
	if !isAuthorizedThroughSoc && !manager.ComparePasswords(user.Password, password) {
		c.JSON(http.StatusUnauthorized, manager.ErrNotCorrectLoginPassword.Error()) // 401
		return
	}

	tokenInfo, session, status, err := serviceSession.ReloadTokens(user.UserID, configs.Cfg.TokenPassword, user.Role, c)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(status, err.Error())
		return
	}

	secure := false
	timeDiff := int(time.Until(session.ExpiresIn) / time.Second)

	if configs.Cfg.Mode == gin.ReleaseMode {
		secure = true
	}

	c.SetCookie(manager.CookieTokenName,
		session.RefreshToken,
		timeDiff, manager.CookiePath, manager.CookieDomain, secure, true)

	userExt := new(response.UserExt)
	userExt.User = user
	userExt.Phones = phones

	c.JSON(200, gin.H{
		"JWT":     tokenInfo.JWT,
		"userExt": userExt,
	})
}

// PostAuthRefreshTokens - обновление токена
func PostAuthRefreshTokens(c *gin.Context) {
	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceSession := service.NewSessionService()
	serviceUsers := service.NewUserService()
	servicePhone := service.NewPhoneService()

	user, err := serviceUsers.GetUserByID(userID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	phones, err := servicePhone.GetPhonesByUserID(userID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tokenInfo, session2, status, err := serviceSession.ReloadTokens(userID, configs.Cfg.TokenPassword, user.Role, c)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(status, err.Error())
		return
	}

	secure := false
	timeDiff := int(time.Until(session2.ExpiresIn) / time.Second)

	if configs.Cfg.Mode == gin.ReleaseMode {
		secure = true
	}

	c.SetCookie(manager.CookieTokenName,
		session2.RefreshToken,
		timeDiff, manager.CookiePath, manager.CookieDomain, secure, true)

	userExt := new(response.UserExt)
	userExt.User = user
	userExt.Phones = phones

	c.JSON(200, gin.H{
		"JWT":     tokenInfo.JWT,
		"userExt": userExt,
	})
}
