package response

import "altair/storage"

// UserExt - структура ответа, пользователь (расширенный)
type UserExt struct {
	*storage.User
	Phones []*storage.Phone `json:"phones"`
}
