package request

type PostCat struct {
	Name     string `form:"name" binding:"required"`
	ParentId uint64 `form:"parentId"`
	Pos      uint64 `form:"pos"`
}
