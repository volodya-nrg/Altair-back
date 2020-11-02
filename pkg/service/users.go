package service

import (
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"gorm.io/gorm"
	"unicode/utf8"
)

// NewUserService - фабрика, создает объект пользователя
func NewUserService() *UserService {
	return new(UserService)
}

// UserService - структура пользователя
type UserService struct{}

// GetUsers - получить пользователей
func (us UserService) GetUsers(order string) ([]*storage.User, error) {
	users := make([]*storage.User, 0)
	err := server.Db.Order(order).Find(&users).Error
	return users, err
}

// GetUserByID - получить данные о пользователе относительно его ID
func (us UserService) GetUserByID(userID uint64) (*storage.User, error) {
	user := new(storage.User)
	err := server.Db.First(user, userID).Error
	return user, err
}

// GetUserByEmail - получить данные о пользователе через e-mail
func (us UserService) GetUserByEmail(userEmail string) (*storage.User, error) {
	user := new(storage.User)
	err := server.Db.Where("email = ?", userEmail).First(user).Error
	return user, err
}

// GetUserByHashCheckEmail - получить данные о пользователе относительно проверочного хеша на е-мэйл
func (us UserService) GetUserByHashCheckEmail(hash string) (*storage.User, error) {
	user := new(storage.User)
	err := server.Db.Where("hash_for_check_email = ?", hash).First(user).Error
	return user, err
}

// HasUser - существует ли пользователь, проверка через валидный е-мэйл
func (us UserService) HasUser(userEmail string) (bool, error) {
	var count int64
	var result bool
	query := `SELECT COUNT(user_id) FROM users WHERE is_email_confirmed = 1 AND email = ?`

	if err := server.Db.Raw(query, userEmail).Count(&count).Error; err != nil {
		return result, err
	}

	result = count > 0

	return result, nil
}

// Create - создать пользователя
func (us UserService) Create(user *storage.User, tx *gorm.DB) error {
	if err := validate(user); err != nil {
		return err
	}
	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(user).Error

	return err
}

// Update - изменить пользователя
func (us UserService) Update(user *storage.User, tx *gorm.DB) error {
	if err := validate(user); err != nil {
		return err
	}
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(user).Error

	return err
}

// Delete - удалить пользователя
func (us UserService) Delete(userID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	if err := tx.Delete(storage.User{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func validate(user *storage.User) error {
	if !manager.ValidateEmail(user.Email) {
		return manager.ErrEmailNotCorrect
	}
	if utf8.RuneCountInString(user.Password) < manager.MinLenPassword {
		return manager.ErrPasswordIsShort
	}

	return nil
}
