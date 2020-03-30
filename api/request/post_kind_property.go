package request

type PostKindProperty struct {
	Name string `form:"name" binding:"required"`
}
