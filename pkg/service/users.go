package service

import (
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
	"unicode/utf8"
)

func NewUserService() *UserService {
	return new(UserService)
}

type UserService struct{}

func (us UserService) Get() bool {
	return true
}
func (us UserService) GetUsers() ([]*storage.User, error) {
	users := make([]*storage.User, 0)
	err := server.Db.Debug().Order("created_at desc").Find(&users).Error
	return users, err
}
func (us UserService) GetUserByID(userId uint64) (*storage.User, error) {
	user := new(storage.User)
	err := server.Db.Debug().First(user, userId).Error
	return user, err
}
func (us UserService) Create(user *storage.User, tx *gorm.DB) error {
	if err := validate(user); err != nil {
		return err
	}
	if !server.Db.Debug().NewRecord(user) {
		return errNotCreateNewUser
	}

	user.Password = helpers.HashAndSalt(user.Password)

	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Create(user).Error

	return err
}
func (us UserService) Update(user *storage.User, tx *gorm.DB) error {
	if err := validate(user); err != nil {
		return err
	}
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(user).Error

	return err
}
func (us UserService) Delete(userId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	if err := tx.Where("user_id = ?", userId).Delete(storage.User{}).Error; err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func validate(user *storage.User) error {
	if !helpers.ValidateEmail(user.Email) {
		return errNotCorrectEmail
	}
	if utf8.RuneCountInString(user.Password) < minLenPassword {
		return errPasswordIsShort
	}

	return nil
}
