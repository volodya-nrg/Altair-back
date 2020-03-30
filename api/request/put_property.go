package request

type PutProperty struct {
	PropertyId     uint64 `form:"propertyId" binding:"required"`
	Title          string `form:"title" binding:"required"`
	KindPropertyId uint64 `form:"kindPropertyId" binding:"required"`
	Name           string `form:"name" binding:"required"`
	MaxInt         uint64 `form:"maxInt"`
	IsRequire      bool   `form:"isRequire"`
	IsCanAsFilter  bool   `form:"isCanAsFilter"`
}
