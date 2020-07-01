package response

import "altair/storage"

// CatTree - структура ответа, дерево каталога
type CatTree struct {
	*storage.Cat
	Childes []*CatTree `json:"childes"`
}
