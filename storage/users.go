package storage

import "time"

type User struct {
	UserId           uint64     `json:"userId" gorm:"primary_key;column:user_id"`
	Email            string     `json:"email" gorm:"column:email"`
	EmailIsConfirmed bool       `json:"emailIsConfirmed" gorm:"column:email_is_confirmed"`
	Name             string     `json:"name" gorm:"column:name"`
	Password         string     `json:"-" gorm:"column:password"`
	Avatar           string     `json:"avatar" gorm:"column:avatar"`
	CreatedAt        *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        *time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
