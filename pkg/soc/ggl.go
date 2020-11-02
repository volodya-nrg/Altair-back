package soc

import (
	"altair/pkg/manager"
	"fmt"
)

// NewGgl - фабрика, создает объект
func NewGgl(clientID, clientSecret, code string) *Ggl {
	output := new(Ggl)
	output.clientID = clientID
	output.clientSecret = clientSecret
	output.code = code

	return output
}

// Ggl - создаем структуру Ggl
type Ggl struct {
	clientID        string
	clientSecret    string
	code            string
	accessTokenData GglAccessToken
}

// GetAccessToken - метод получающий токен
func (s *Ggl) GetAccessToken(redirect string) error {
	url := "https://oauth2.googleapis.com/token"
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
func (s *Ggl) GetUserInfo() (interface{}, error) {
	userInfo := GglUserInfo{}
	query := map[string]string{
		"access_token": s.accessTokenData.AccessToken,
	}

	if err := s.checkAccessToken(); err != nil {
		return userInfo, err
	}

	if err := manager.MakeRequest("get", "https://www.googleapis.com/oauth2/v1/userinfo", &userInfo, query); err != nil {
		return userInfo, err
	}

	return userInfo, nil
}
func (s *Ggl) checkAccessToken() error {
	if s.accessTokenData.AccessToken == "" || s.accessTokenData.ExpiresIn < 1 {
		return manager.ErrAccessTokenNotCorrect
	}

	return nil
}

// GglAccessToken - структура ответа AccessToken
type GglAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// GglUserInfo - структура ответа о пользователе
type GglUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}
