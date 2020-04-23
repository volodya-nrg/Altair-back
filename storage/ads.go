package storage

import "time"

type Ad struct {
	AdId        uint64     `json:"adId" gorm:"primary_key;column:ad_id"`
	Title       string     `json:"title" gorm:"column:title"`
	Slug        string     `json:"slug" gorm:"column:slug"`
	CatId       uint64     `json:"catId" gorm:"column:cat_id"`
	UserId      uint64     `json:"userId" gorm:"column:user_id"`
	Description string     `json:"description" gorm:"column:description"`
	Price       uint64     `json:"price" gorm:"column:price"`
	Ip          string     `json:"ip" gorm:"column:ip"`
	IsDisabled  bool       `json:"isDisabled" gorm:"column:is_disabled"`
	HasPhoto    bool       `json:"hasPhoto" gorm:"column:has_photo"`
	Youtube     string     `json:"youtube" gorm:"column:youtube"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
