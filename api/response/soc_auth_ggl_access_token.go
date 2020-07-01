package response

// SocAuthGglAccessToken - структура ответа Гугла, получение токена
type SocAuthGglAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
