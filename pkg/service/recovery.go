package service

import (
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

// NewRecoveryService - фабрика, создает объект "восстановление"
func NewRecoveryService() *RecoveryService {
	return new(RecoveryService)
}

// RecoveryService - структура "восстановления"
type RecoveryService struct{}

// GetByHash - получить данные относительно хеша
func (rs RecoveryService) GetByHash(hash string) (*storage.Recovery, error) {
	hashRow := new(storage.Recovery)

	err := server.Db.Order("created_at desc").Where("hash = ?", hash).First(hashRow).Error

	return hashRow, err
}

// Create - создать запись
func (rs RecoveryService) Create(recovery *storage.Recovery, tx *gorm.DB) error {
	if !server.Db.NewRecord(recovery) {
		return manager.ErrNotCreateNewRecovery
	}

	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(recovery).Error

	return err
}

// Delete - удалить запись
func (rs RecoveryService) Delete(userID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	if err := tx.Delete(storage.Recovery{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}
