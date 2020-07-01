package request

// PutUser - структура запроса на измение пользователя
type PutUser struct {
	UserID           uint64 `form:"userId" binding:"required"`
	Email            string `form:"email" binding:"required"`
	Avatar           string `form:"avatar"`
	Name             string `form:"name"`
	Password         string `form:"password"`
	PasswordConfirm  string `form:"passwordConfirm"`
	IsEmailConfirmed bool   `form:"isEmailConfirmed"`
}
