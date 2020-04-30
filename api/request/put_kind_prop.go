package request

type PutKindProp struct {
	KindPropId uint64 `form:"kindPropId" binding:"required"`
	Name       string `form:"name" binding:"required"`
}
