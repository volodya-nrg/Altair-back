package request

type GetSearchAds struct {
	Query  string `form:"q" binding:"required"`
	CatId  uint64 `form:"catId"`
	Limit  uint64 `form:"limit"`
	Offset uint64 `form:"offset"`
}
