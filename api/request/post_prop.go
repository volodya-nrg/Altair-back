package request

type PostProp struct {
	Title            string `form:"title" binding:"required"`
	KindPropId       uint64 `form:"kindPropId" binding:"required"`
	Name             string `form:"name" binding:"required"`
	Suffix           string `form:"suffix"`
	Comment          string `form:"comment"`
	PrivateComment   string `form:"privateComment"`
	SelectAsTextarea string `form:"select_as_textarea"`
}
