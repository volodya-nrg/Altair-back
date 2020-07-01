package response

// CatTreeFull - структура ответа, дерево каталога (полное)
type CatTreeFull struct {
	*СatFull
	Childes []*CatTreeFull `json:"childes"`
}
