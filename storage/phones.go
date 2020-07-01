package storage

import "time"

// Phone - структура таблицы номеров телефонов
type Phone struct {
	PhoneID   uint64    `json:"phoneId" gorm:"primary_key;column:phone_id"`
	Number    string    `json:"number" gorm:"column:number"`
	IsVerify  bool      `json:"isVerify" gorm:"column:is_verify"`
	UserID    uint64    `json:"-" gorm:"column:user_id"`
	Code      string    `json:"-" gorm:"column:code"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
