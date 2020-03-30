package request

type PutKindProperty struct {
	KindPropertyId uint64 `form:"kindPropertyId" binding:"required"`
	Name           string `form:"name" binding:"required"`
}
