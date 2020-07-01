package response

// SocAuthVkAccessToken - структура ответа VK, получение токена
type SocAuthVkAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
}
