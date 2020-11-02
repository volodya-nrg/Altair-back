package soc

import (
	"altair/pkg/manager"
	"fmt"
)

// NewFb - фабрика, создает объект
func NewFb(clientID uint64, clientSecret, code string) *Fb {
	output := new(Fb)
	output.clientID = clientID
	output.clientSecret = clientSecret
	output.code = code

	return output
}

// Fb - создаем структуру Fb
type Fb struct {
	clientID        uint64
	clientSecret    string
	code            string
	accessTokenData FbAccessToken
}

// GetAccessToken - метод получающий токен
func (s *Fb) GetAccessToken(redirect string) error {
	url := "https://graph.facebook.com/v7.0/oauth/access_token"
	query := map[string]string{
		"client_id":     fmt.Sprint(s.clientID),
		"client_secret": s.clientSecret,
		"redirect_uri":  redirect,
		"code":          s.code,
	}

	if err := manager.MakeRequest("post", url, &s.accessTokenData, query); err != nil {
		return err
	}

	return nil
}

// GetUserInfo - метод получающий информацию о пользователе
func (s *Fb) GetUserInfo() (interface{}, error) {
	userInfo := FbUserInfo{}
	query := map[string]string{
		"access_token": s.accessTokenData.AccessToken,
	}

	if err := s.checkAccessToken(); err != nil {
		return userInfo, err
	}

	if err := manager.MakeRequest("get", "https://graph.facebook.com/me", &userInfo, query); err != nil {
		return userInfo, err
	}

	return userInfo, nil
}
func (s *Fb) checkAccessToken() error {
	if s.accessTokenData.AccessToken == "" || s.accessTokenData.ExpiresIn < 1 {
		return manager.ErrAccessTokenNotCorrect
	}

	return nil
}

// FbAccessToken - структура ответа AccessToken
type FbAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// FbUserInfo - структура ответа о пользователе
type FbUserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
