package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/configs"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/storage"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
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
		var tempEmail string

		if state == "vk" {
			responseVkAccessToken := new(response.SocAuthVkAccessToken)
			query := map[string]string{
				"client_id":     fmt.Sprint(configs.Cfg.Socials.Vk.ClientID),
				"client_secret": configs.Cfg.Socials.Vk.ClientSecret,
				"redirect_uri":  "https://www.altair.uz/login", // так же как и на фронте
				"code":          code,
			}

			if err := manager.MakeRequest("post", "https://oauth.vk.com/access_token", responseVkAccessToken, query); err != nil {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return
			}

			if responseVkAccessToken.AccessToken != "" && responseVkAccessToken.ExpiresIn > 0 && responseVkAccessToken.UserID > 0 {
				tempEmail = fmt.Sprintf("id%d@vk.com", responseVkAccessToken.UserID)

			} else {
				c.JSON(400, manager.ErrSocAuthUnknown.Error())
				return
			}
		} else if state == "ok" {
			responseOkAccessToken := new(response.SocAuthOkAccessToken)
			responseOkCurrentUser := new(response.SocAuthOkCurrentUser)
			query := map[string]string{
				"code":          code,
				"client_id":     fmt.Sprint(configs.Cfg.Socials.Ok.ClientID),
				"client_secret": configs.Cfg.Socials.Ok.ClientSecret,
				"redirect_uri":  "https://www.altair.uz/login", // так же как и на фронте
				"grant_type":    "authorization_code",
			}

			if err := manager.MakeRequest("post", "https://api.ok.ru/oauth/token.do", responseOkAccessToken, query); err != nil {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return
			}

			if responseOkAccessToken.AccessToken != "" && responseOkAccessToken.ExpiresIn != "" {
				// получаем данные о пользователе
				str1 := fmt.Sprintf("%s%s", responseOkAccessToken.AccessToken, configs.Cfg.Socials.Ok.ClientSecret)
				md5Hash1 := fmt.Sprintf("%x", md5.Sum([]byte(str1)))
				preSig := []string{"application_key=" + configs.Cfg.Socials.Ok.ClientPublic, "fields=uid", "format=json", "method=users.getCurrentUser", md5Hash1}
				str2 := strings.Join(preSig, "")
				md5Hash2 := fmt.Sprintf("%x", md5.Sum([]byte(str2)))
				query2 := map[string]string{
					"application_key": configs.Cfg.Socials.Ok.ClientPublic,
					"fields":          "uid",
					"format":          "json",
					"method":          "users.getCurrentUser",
					"sig":             md5Hash2,
					"access_token":    responseOkAccessToken.AccessToken,
				}

				if err := manager.MakeRequest("post", "https://api.ok.ru/fb.do", responseOkCurrentUser, query2); err != nil {
					logger.Warning.Println(err.Error())
					c.JSON(500, err.Error())
					return
				}

				if responseOkCurrentUser.UID != "" {
					tempEmail = fmt.Sprintf("id%s@ok.ru", responseOkCurrentUser.UID)
				}

			} else {
				c.JSON(400, manager.ErrSocAuthUnknown.Error())
				return
			}
		} else if state == "fb" {
			responseFbAccessToken := new(response.SocAuthFbAccessToken)
			responseFbCurrentUser := new(response.SocAuthFbCurrentUser)
			query := map[string]string{
				"client_id":     fmt.Sprint(configs.Cfg.Socials.Fb.ClientID),
				"redirect_uri":  "https://www.altair.uz/login", // так же как и на фронте
				"client_secret": configs.Cfg.Socials.Fb.ClientSecret,
				"code":          code,
			}
			if err := manager.MakeRequest("post", "https://graph.facebook.com/v7.0/oauth/access_token", responseFbAccessToken, query); err != nil {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return
			}

			if responseFbAccessToken.AccessToken != "" && responseFbAccessToken.ExpiresIn > 0 {
				query := map[string]string{
					"access_token": responseFbAccessToken.AccessToken,
				}

				if err := manager.MakeRequest("get", "https://graph.facebook.com/me", responseFbCurrentUser, query); err != nil {
					logger.Warning.Println(err.Error())
					c.JSON(500, err.Error())
					return
				}

				if responseFbCurrentUser.ID != "" {
					tempEmail = fmt.Sprintf("id%s@facebook.com", responseFbCurrentUser.ID)
				}

			} else {
				c.JSON(400, manager.ErrSocAuthUnknown.Error())
				return
			}
		} else if state == "ggl" {
			responseGglAccessToken := new(response.SocAuthGglAccessToken)
			query := map[string]string{
				"client_id":     configs.Cfg.Socials.Ggl.ClientID,
				"client_secret": configs.Cfg.Socials.Ggl.ClientSecret,
				"code":          code,
				"grant_type":    "authorization_code",
				"redirect_uri":  "https://www.altair.uz/login", // так же как и на фронте
			}
			if err := manager.MakeRequest("post", "https://oauth2.googleapis.com/token", responseGglAccessToken, query); err != nil {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return
			}

			logger.Info.Printf("%#v", responseGglAccessToken)

			if responseGglAccessToken.AccessToken != "" && responseGglAccessToken.ExpiresIn > 0 {
				var x interface{}
				query := map[string]string{
					"access_token": responseGglAccessToken.AccessToken,
				}
				if err := manager.MakeRequest("post", "https://www.googleapis.com/oauth2/v1/userinfo", x, query); err != nil {
					logger.Warning.Println(err.Error())
					c.JSON(500, err.Error())
					return
				}

				logger.Info.Printf("%#v", x)

			} else {
				c.JSON(400, manager.ErrSocAuthUnknown.Error())
				return
			}
		}

		if tempEmail != "" {
			email = tempEmail
			isAuthorizedThroughSoc = true

			_, err := serviceUsers.GetUserByEmail(tempEmail)
			if err != nil && !gorm.IsRecordNotFoundError(err) {
				logger.Warning.Println(err.Error())
				c.JSON(500, err.Error())
				return

			} else if gorm.IsRecordNotFoundError(err) {
				tmpUser := new(storage.User)
				tmpUser.Email = tempEmail
				tmpUser.Password = manager.HashAndSalt(manager.RandStringRunes(10))
				tmpUser.IsEmailConfirmed = true

				if err := serviceUsers.Create(tmpUser, nil); err != nil {
					logger.Warning.Println(err.Error())
					c.JSON(500, err.Error())
					return
				}
			}
		}
	}

	// стандартный подход проверки пользователя
	user, err := serviceUsers.GetUserByEmail(email)
	if gorm.IsRecordNotFoundError(err) {
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
