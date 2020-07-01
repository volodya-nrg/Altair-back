package request

// PostRecoverySendHash - структура запроса на добавление записи о восстановление е-мэйла
type PostRecoverySendHash struct {
	Email string `form:"email" binding:"required"`
}
