package request

// PostKindProp - структура запроса на добавление вида свойства
type PostKindProp struct {
	Name string `form:"name" binding:"required"`
}
