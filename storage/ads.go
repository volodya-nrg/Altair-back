package storage

import "time"

// Ad - структура таблицы объявлений
type Ad struct {
	AdID        uint64     `json:"adID" gorm:"primaryKey;column:ad_id"`
	Title       string     `json:"title" gorm:"column:title"`
	Slug        string     `json:"slug" gorm:"column:slug"`
	CatID       uint64     `json:"catID" gorm:"column:cat_id"`
	UserID      uint64     `json:"userID" gorm:"column:user_id"`
	Description string     `json:"description" gorm:"column:description"`
	Price       uint64     `json:"price" gorm:"column:price"`
	IP          string     `json:"IP" gorm:"column:ip"`
	IsDisabled  bool       `json:"isDisabled" gorm:"column:is_disabled"`
	IsApproved  bool       `json:"isApproved" gorm:"column:is_approved"`
	HasPhoto    bool       `json:"hasPhoto" gorm:"column:has_photo"`
	Youtube     string     `json:"youtube" gorm:"column:youtube"`
	Latitude    float64    `json:"latitude" gorm:"column:latitude"`
	Longitude   float64    `json:"longitude" gorm:"column:longitude"`
	CityName    string     `json:"cityName" gorm:"column:city_name"`
	CountryName string     `json:"countryName" gorm:"column:country_name"`
	PhoneID     uint64     `json:"phoneID" gorm:"column:phone_id"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
