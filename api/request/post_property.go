package request

type PostProperty struct {
	Title          string `form:"title" binding:"required"`
	KindPropertyId uint64 `form:"kindPropertyId" binding:"required"`
	Name           string `form:"name" binding:"required"`
	IsRequire      bool   `form:"isRequire"`
}
