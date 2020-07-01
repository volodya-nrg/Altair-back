package response

import "altair/storage"

// СatFull - структура ответа, категория (полное)
type СatFull struct {
	*storage.Cat
	PropsFull []*PropFull `json:"props"`
}
