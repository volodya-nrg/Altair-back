package request

// PostProfile - структура запроса на добавление профиля
type PostProfile struct {
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"passwordConfirm" binding:"required"`
	AgreeOffer      bool   `form:"agreeOffer" binding:"required"`
	AgreePolicy     bool   `form:"agreePolicy" binding:"required"`
}
