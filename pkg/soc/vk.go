package soc

import (
	"altair/pkg/manager"
	"fmt"
)

// NewVk - фабрика, создает объект
func NewVk(clientID uint64, clientSecret, code string) *Vk {
	output := new(Vk)
	output.clientID = clientID
	output.clientSecret = clientSecret
	output.code = code

	return output
}

// Vk - структура VK
type Vk struct {
	clientID        uint64
	clientSecret    string
	code            string
	accessTokenData VkAccessToken
}

// GetAccessToken - метод получающий токен
func (s *Vk) GetAccessToken(redirect string) error {
	url := "https://oauth.vk.com/access_token"
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
func (s *Vk) GetUserInfo() (interface{}, error) {
	response := new(VkResponseUserInfo)
	userInfo := VkUserInfo{}
	url := "https://api.vk.com/method/users.get"
	query := map[string]string{
		"access_token": s.accessTokenData.AccessToken,
		"v":            "5.120",
	}

	if err := s.checkAccessToken(); err != nil {
		return userInfo, err
	}

	if err := manager.MakeRequest("post", url, response, query); err != nil {
		return userInfo, err
	}

	userInfo = response.Response[0]

	return userInfo, nil
}
func (s *Vk) checkAccessToken() error {
	if s.accessTokenData.AccessToken == "" || s.accessTokenData.ExpiresIn < 1 || s.accessTokenData.UserID < 1 {
		return manager.ErrAccessTokenNotCorrect
	}

	return nil
}

// VkAccessToken - структура ответа AccessToken
type VkAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
}

// VkResponseUserInfo - структура ответа (общая)
type VkResponseUserInfo struct {
	Response []VkUserInfo `json:"response"`
}

// VkUserInfo - структура ответа о пользователе
type VkUserInfo struct {
	ID              uint64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	IsClosed        bool   `json:"is_closed"`         // включена ли приватность профиля
	CanAccessClosed bool   `json:"can_access_closed"` // есть ли у текущего пользователя возможность видеть профиль пользователя при is_closed = true
}
