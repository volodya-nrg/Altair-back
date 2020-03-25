package service

import (
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"errors"
	"unicode/utf8"
)

var (
	errNotCorrectEmail  = errors.New("not correct email")
	errPasswordIsShort  = errors.New("password is short")
	errNotCreateNewUser = errors.New("not create new user")
	minLenPassword      = 6
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
	return users, server.Db.Debug().Order("user_id", true).Find(users).Error
}
func (userService UserService) GetUserByID(userId uint64) (*storage.User, error) {
	user := new(storage.User)
	return user, server.Db.Debug().First(user, userId).Error // проверяется в контроллере
}
func (userService UserService) Create(user *storage.User) error {
	if err := validate(user); err != nil {
		return err
	}
	if !server.Db.Debug().NewRecord(user) {
		return errNotCreateNewUser
	}

	user.Password = helpers.HashAndSalt(user.Password)

	return server.Db.Debug().Create(user).Error
}
func (userService UserService) Update(user *storage.User) error {
	if err := validate(user); err != nil {
		return err
	}

	return server.Db.Debug().Save(user).Error
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
