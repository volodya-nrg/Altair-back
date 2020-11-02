package request

// PutAd - структура запроса на изменение объявления
type PutAd struct {
	PostAd
	AdID uint64 `form:"adID" binding:"required"`
}
