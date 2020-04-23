package request

type PostProperty struct {
	Title            string `form:"title" binding:"required"`
	KindPropertyId   uint64 `form:"kindPropertyId" binding:"required"`
	Name             string `form:"name" binding:"required"`
	Suffix           string `form:"suffix"`
	Comment          string `form:"comment"`
	PrivateComment   string `form:"privateComment"`
	SelectAsTextarea string `form:"select_as_textarea"`
}
