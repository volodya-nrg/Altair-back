package request

import "altair/storage"

// PostProp - структура запроса на добавление свойства
type PostProp struct {
	Title            string              `form:"title" binding:"required"`
	KindPropID       string              `form:"kindPropID" binding:"required"`
	Name             string              `form:"name" binding:"required"`
	Suffix           string              `form:"suffix"`
	Comment          string              `form:"comment"`
	PrivateComment   string              `form:"privateComment"`
	SelectAsTextarea string              `form:"selectAsTextarea"`
	Values           []storage.ValueProp `form:"values"`
}
