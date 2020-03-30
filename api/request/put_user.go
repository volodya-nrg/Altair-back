package request

type PutUser struct {
	UserId           uint64 `form:"userId" binding:"required"`
	Email            string `form:"email" binding:"required"`
	Avatar           string `form:"avatar"`
	Name             string `form:"name"`
	PasswordOld      string `form:"passwordOld"`
	Password         string `form:"password"`
	PasswordConfirm  string `form:"passwordConfirm"`
	EmailIsConfirmed bool   `form:"emailIsConfirmed"`
	// File             *multipart.FileHeader `form:"file"`
}
