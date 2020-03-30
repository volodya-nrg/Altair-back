package service

import (
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"unicode/utf8"
)

func NewUserService() *UserService {
	return new(UserService)
}

type UserService struct{}

func (userService UserService) Get() bool {
	return true
}
func (userService UserService) GetUsers() ([]*storage.User, error) {
	users := make([]*storage.User, 0)
	err := server.Db.Debug().Order("user_id", true).Find(&users).Error
	return users, err
}
func (userService UserService) GetUserByID(userId uint64) (*storage.User, error) {
	user := new(storage.User)
	err := server.Db.Debug().First(user, userId).Error // проверяется в контроллере
	return user, err
}
func (userService UserService) Create(user *storage.User) error {
	if err := validate(user); err != nil {
		return err
	}
	if !server.Db.Debug().NewRecord(user) {
		return errNotCreateNewUser
	}

	user.Password = helpers.HashAndSalt(user.Password)
	err := server.Db.Debug().Create(user).Error

	return err
}
func (userService UserService) Update(user *storage.User) error {
	if err := validate(user); err != nil {
		return err
	}

	err := server.Db.Debug().Save(user).Error
	return err
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
