package storage

import "time"

// Recovery - структура таблицы для восстановления пароля
type Recovery struct {
	RecoverID uint64     `json:"recoverID" gorm:"primaryKey;column:recover_id"`
	Hash      string     `json:"hash" gorm:"column:hash"`
	UserID    uint64     `json:"userID" gorm:"column:user_id"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at"`
}
