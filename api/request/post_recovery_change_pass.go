package request

// PostRecoveryChangePass - структура запроса на добавление записи связанным с изменение пароля
type PostRecoveryChangePass struct {
	Hash            string `form:"hash" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"passwordConfirm" binding:"required"`
}
