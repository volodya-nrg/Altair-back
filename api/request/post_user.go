package request

type PostUser struct {
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"passwordConfirm" binding:"required"`
	AgreeOffer      bool   `form:"agreeOffer" binding:"required"`
	AgreePolicy     bool   `form:"agreePolicy" binding:"required"`
}
