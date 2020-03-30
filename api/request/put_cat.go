package request

type PutCat struct {
	CatId      uint64 `form:"catId" binding:"required"`
	Name       string `form:"name" binding:"required"`
	ParentId   uint64 `form:"parentId"`
	Pos        uint64 `form:"pos"`
	IsDisabled bool   `form:"isDisabled"`
}
