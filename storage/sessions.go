package storage

import "time"

// Session - структура таблицы сессий
type Session struct {
	SessionID    uint64    `json:"sessionID" gorm:"primaryKey;column:session_id"`
	UserID       uint64    `json:"userID" gorm:"column:user_id"`
	RefreshToken string    `json:"refreshToken" gorm:"column:refresh_token"`
	ExpiresIn    time.Time `json:"expiresIn" gorm:"column:expires_in"`
	UserAgent    string    `json:"userAgent" gorm:"column:user_agent"`
	IP           string    `json:"IP" gorm:"column:ip"`
	Fingerprint  string    `json:"fingerprint" gorm:"column:fingerprint"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
}
