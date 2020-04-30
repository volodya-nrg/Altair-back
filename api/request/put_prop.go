package request

type PutProp struct {
	PropId         uint64 `form:"propId" binding:"required"`
	Title          string `form:"title" binding:"required"`
	KindPropId     uint64 `form:"kindPropId" binding:"required"`
	Name           string `form:"name" binding:"required"`
	IsRequire      bool   `form:"isRequire"`
	Suffix         string `form:"suffix"`
	Comment        string `form:"comment"`
	PrivateComment string `form:"privateComment"`
}
