package request

// PostUser - структура запроса на добавление пользователя
type PostUser struct {
	Email            string `form:"email" binding:"required"`
	Name             string `form:"name"`
	Password         string `form:"password" binding:"required"`
	PasswordConfirm  string `form:"passwordConfirm" binding:"required"`
	IsEmailConfirmed bool   `form:"isEmailConfirmed"`
}
