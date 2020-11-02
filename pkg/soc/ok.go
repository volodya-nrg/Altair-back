package soc

import (
	"altair/configs"
	"altair/pkg/manager"
	"crypto/md5"
	"fmt"
)

// NewOk - фабрика, создает объект
func NewOk(clientID uint64, clientSecret, clientPublic, code string) *Ok {
	output := new(Ok)
	output.clientID = clientID
	output.clientSecret = clientSecret
	output.clientPublic = clientPublic
	output.code = code

	return output
}

// Ok - создаем структуру OK
type Ok struct {
	clientID        uint64
	clientSecret    string
	clientPublic    string
	code            string
	accessTokenData OkAccessToken
}

// GetAccessToken - метод получающий токен
func (s *Ok) GetAccessToken(redirect string) error {
	url := "https://api.ok.ru/oauth/token.do"
	query := map[string]string{
		"client_id":     fmt.Sprint(s.clientID),
		"client_secret": s.clientSecret,
		"redirect_uri":  redirect,
		"code":          s.code,
		"grant_type":    "authorization_code",
	}

	if err := manager.MakeRequest("post", url, &s.accessTokenData, query); err != nil {
		return err
	}

	return nil
}

// GetUserInfo - метод получающий информацию о пользователе
func (s *Ok) GetUserInfo() (interface{}, error) {
	userInfo := OkUserInfo{}
	str1 := fmt.Sprintf("%s%s", s.accessTokenData.AccessToken, s.clientSecret)
	str2 := fmt.Sprintf("application_key=%sfields=uid,first_name,last_nameformat=jsonmethod=users.getCurrentUser%x", s.clientPublic, md5.Sum([]byte(str1)))
	query := map[string]string{
		"application_key": configs.Cfg.Socials.Ok.ClientPublic,
		"fields":          "uid,first_name,last_name",
		"format":          "json",
		"method":          "users.getCurrentUser",
		"sig":             fmt.Sprintf("%x", md5.Sum([]byte(str2))),
		"access_token":    s.accessTokenData.AccessToken,
	}

	if err := s.checkAccessToken(); err != nil {
		return userInfo, err
	}

	if err := manager.MakeRequest("post", "https://api.ok.ru/fb.do", &userInfo, query); err != nil {
		return userInfo, err
	}

	return userInfo, nil
}
func (s *Ok) checkAccessToken() error {
	if s.accessTokenData.AccessToken == "" || s.accessTokenData.ExpiresIn == "" {
		return manager.ErrAccessTokenNotCorrect
	}

	return nil
}

// OkAccessToken - структура ответа AccessToken
type OkAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

// OkUserInfo - структура ответа о пользователе
type OkUserInfo struct {
	UID       string `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
