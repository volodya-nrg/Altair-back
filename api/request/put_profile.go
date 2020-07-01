package request

// PutProfile - структура запроса изменение профиля
type PutProfile struct {
	Avatar          string `form:"avatar"`
	Name            string `form:"name"`
	PasswordOld     string `form:"passwordOld"`
	PasswordNew     string `form:"passwordNew"`
	PasswordConfirm string `form:"passwordConfirm"`
}
