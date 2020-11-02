package request

// PutKindProp - структура запроса изменение вида свойства
type PutKindProp struct {
	KindPropID uint64 `form:"kindPropID" binding:"required"`
	Name       string `form:"name" binding:"required"`
}
