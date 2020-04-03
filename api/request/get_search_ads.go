package request

type GetSearchAds struct {
	Query string `form:"q" binding:"required"`
	CatId uint64 `form:"catId"`
}
