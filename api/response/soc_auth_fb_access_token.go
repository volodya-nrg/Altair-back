package response

// SocAuthFbAccessToken - структура ответа Фейсбука, получение токета
type SocAuthFbAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SocAuthFbCurrentUser - структура ответа Фейсбука, получение данных пользователя
type SocAuthFbCurrentUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
