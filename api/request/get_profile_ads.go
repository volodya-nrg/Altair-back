package request

// GetProfileAds - структура запроса на объявления пользователя
type GetProfileAds struct {
	Limit  uint64 `form:"limit"`
	Offset uint64 `form:"offset"`
}
