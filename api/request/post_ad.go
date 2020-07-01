package request

// PostAd - структура запроса на добавление объявления
type PostAd struct {
	Title       string  `form:"title"` // может быть и не обязательным. Нужно помнить про Slug
	CatID       uint64  `form:"catId" binding:"required"`
	Description string  `form:"description" binding:"required"`
	PhoneID     uint64  `form:"phoneId" binding:"required"`
	Price       uint64  `form:"price"` // может быть ноль
	UserID      uint64  `form:"userId"`
	Youtube     string  `form:"youtube"`
	Latitude    float64 `form:"latitude"`
	Longitude   float64 `form:"longitude"`
	CityName    string  `form:"cityName"`
	CountryName string  `form:"countryName"`
	IsDisabled  bool    `form:"isDisabled"`
	IsApproved  bool    `form:"isApproved"`
}
