package request

// PutKindProp - структура запроса изменение вида свойства
type PutKindProp struct {
	KindPropID uint64 `form:"kindPropId" binding:"required"`
	Name       string `form:"name" binding:"required"`
}
