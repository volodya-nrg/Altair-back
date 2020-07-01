package request

// PostProfilePhone - структура запроса на добавление номера телефона к пользователю
type PostProfilePhone struct {
	Number string `form:"number" binding:"required"`
}
