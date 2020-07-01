package service

import (
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

// NewPhoneService - фабрика, создает объект номера телефона
func NewPhoneService() *PhoneService {
	return new(PhoneService)
}

// PhoneService - структура о номере телефона
type PhoneService struct{}

// GetByID - получить относительно от его ID
func (ps PhoneService) GetByID(phoneID uint64) (*storage.Phone, error) {
	phone := new(storage.Phone)
	err := server.Db.Where("is_verify = ?", true).First(phone, phoneID).Error
	return phone, err
}

// GetPhoneByNumberAndUserID - получить номер телефона относительно его номера и ID пользователя
func (ps PhoneService) GetPhoneByNumberAndUserID(number string, userID uint64) (*storage.Phone, error) {
	phone := new(storage.Phone)
	err := server.Db.Where("number = ? AND user_id = ?", number, userID).First(phone).Error
	return phone, err
}

// GetPhonesByUserID - получить номера телефонов относительно ID пользователя
func (ps PhoneService) GetPhonesByUserID(userID uint64) ([]*storage.Phone, error) {
	phones := make([]*storage.Phone, 0)
	err := server.Db.Where("user_id = ? AND is_verify = ?", userID, true).Find(&phones).Error
	return phones, err
}

// Create - создать запись
func (ps PhoneService) Create(phone *storage.Phone, tx *gorm.DB) error {
	if !server.Db.NewRecord(phone) {
		return manager.ErrNotCreateNewPhone
	}

	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(phone).Error

	return err
}

// Update - изменить запись
func (ps PhoneService) Update(phone *storage.Phone, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(phone).Error

	return err
}

// Delete - удалить запись
func (ps PhoneService) Delete(number string, userID uint64, onlyNotVerify bool, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	query := tx.Where("number = ? AND user_id = ?", number, userID)

	if onlyNotVerify {
		query = tx.Where("is_verify = ?", 0)
	}

	if err := query.Delete(storage.Phone{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAllByUserID - удалить все записи относительно ID пользователя
func (ps PhoneService) DeleteAllByUserID(userID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	if err := tx.Delete(storage.Phone{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
