package request

type PostKindProp struct {
	Name string `form:"name" binding:"required"`
}
