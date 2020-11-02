package storage

import "time"

// User - структура таблицы пользователей
type User struct {
	UserID            uint64     `json:"userID" gorm:"primaryKey;column:user_id"`
	Email             string     `json:"email" gorm:"column:email"`
	IsEmailConfirmed  bool       `json:"isEmailConfirmed" gorm:"column:is_email_confirmed"`
	Name              string     `json:"name" gorm:"column:name"`
	Password          string     `json:"-" gorm:"column:password"`
	Avatar            string     `json:"avatar" gorm:"column:avatar"`
	Role              string     `json:"-" gorm:"column:role"`
	HashForCheckEmail string     `json:"-" gorm:"column:hash_for_check_email"`
	CreatedAt         *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         *time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
