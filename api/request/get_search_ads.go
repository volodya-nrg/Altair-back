package request

// GetSearchAds - структура запроса на получение объявлений через поиск
type GetSearchAds struct {
	Query  string `form:"q" binding:"required"`
	CatID  uint64 `form:"catId"`
	Limit  uint64 `form:"limit"`
	Offset uint64 `form:"offset"`
}
