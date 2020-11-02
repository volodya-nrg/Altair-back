package request

import "altair/storage"

// PutProp - структура запроса на изменение свойства
type PutProp struct {
	PropID         uint64              `form:"propID" binding:"required"`
	Title          string              `form:"title" binding:"required"`
	KindPropID     string              `form:"kindPropID" binding:"required"`
	Name           string              `form:"name" binding:"required"`
	IsRequire      bool                `form:"isRequire"`
	Suffix         string              `form:"suffix"`
	Comment        string              `form:"comment"`
	PrivateComment string              `form:"privateComment"`
	Values         []storage.ValueProp `form:"values"`
}
