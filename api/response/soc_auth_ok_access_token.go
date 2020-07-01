package response

// SocAuthOkAccessToken - структура ответа OK, получение токена
type SocAuthOkAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

// SocAuthOkCurrentUser - структура ответа OK, получение данных пользователя
type SocAuthOkCurrentUser struct {
	UID string `json:"uid"`
}
