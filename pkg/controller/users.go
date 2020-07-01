package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
	"unicode/utf8"
)

// GetUsers - получение всех пользователей
func GetUsers(c *gin.Context) {
	serviceUsers := service.NewUserService()

	users, err := serviceUsers.GetUsers("created_at desc")
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, users)
}

// GetUsersUserID - получение конкретного пользователя
func GetUsersUserID(c *gin.Context) {
	sUserID := c.Param("userID")
	serviceUsers := service.NewUserService()

	userID, err := manager.SToUint64(sUserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	user, err := serviceUsers.GetUserByID(userID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, user)
}

// PostUsers - добавление пользователя
func PostUsers(c *gin.Context) {
	postRequest := new(request.PostUser)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceUsers := service.NewUserService()

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	password := strings.TrimSpace(postRequest.Password)
	passwordConfirm := strings.TrimSpace(postRequest.PasswordConfirm)

	if utf8.RuneCountInString(password) < manager.MinLenPassword {
		if password != passwordConfirm {
			c.JSON(400, manager.ErrPasswordIsShort.Error())
			return
		}
	}
	if password != passwordConfirm {
		c.JSON(400, manager.ErrPasswordsAreNotEqual.Error())
		return
	}

	// тут надо проверить уникальность пользователя
	if has, err := serviceUsers.HasUser(postRequest.Email); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return

	} else if has {
		c.JSON(400, manager.ErrUserAlreadyExists.Error())
		return
	}

	var filePath string
	if len(form.File["files"]) > 0 {
		file := form.File["files"][0] // только один файл
		filePath, err = manager.UploadImage(file, manager.DirImages, c.SaveUploadedFile)
		if err != nil {
			logger.Warning.Println(err.Error())
		}
	}

	user := new(storage.User)

	user.Name = strings.TrimSpace(postRequest.Name)
	user.Email = strings.TrimSpace(postRequest.Email)
	user.Password = manager.HashAndSalt(password)
	user.Avatar = filePath
	user.IsEmailConfirmed = postRequest.IsEmailConfirmed

	if err := serviceUsers.Create(user, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(201, user)
}

// PutUsersUserID - редактирование пользователя
func PutUsersUserID(c *gin.Context) {
	sUserID := c.Param("userID")
	putRequest := new(request.PutUser)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	userID, err := manager.SToUint64(sUserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	serviceUsers := service.NewUserService()

	user, err := serviceUsers.GetUserByID(userID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	user.Name = strings.TrimSpace(putRequest.Name)

	if user.IsEmailConfirmed != putRequest.IsEmailConfirmed {
		user.IsEmailConfirmed = !user.IsEmailConfirmed
	}

	password := strings.TrimSpace(putRequest.Password)
	passwordConfirm := strings.TrimSpace(putRequest.PasswordConfirm)

	if utf8.RuneCountInString(password) > 0 && utf8.RuneCountInString(passwordConfirm) > 0 {
		if utf8.RuneCountInString(password) < manager.MinLenPassword {
			c.JSON(400, manager.ErrPasswordIsShort.Error())
			return
		}
		if password != passwordConfirm {
			c.JSON(400, manager.ErrPasswordsAreNotEqual.Error())
			return
		}
		user.Password = manager.HashAndSalt(password)
	}

	var filePath string
	if len(form.File["files"]) > 0 {
		file := form.File["files"][0] // только один файл
		filePath, err = manager.UploadImage(file, manager.DirImages, c.SaveUploadedFile)
		if err != nil {
			logger.Warning.Println(err.Error())
		}
	}

	// отправили файл или у Юзера есть аватар и пришедший аватар пустой (удалим)
	if filePath != "" || (user.Avatar != "" && putRequest.Avatar == "") {
		if user.Avatar != "" {
			// файл хранится на удаленном сервере. Тут если делать, то удалять с удаленного сервера.
			user.Avatar = ""
		}
		if filePath != "" {
			user.Avatar = filePath
		}
	}

	if err = serviceUsers.Update(user, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, user)
}

// DeleteUsersUserID - удаление пользователя
func DeleteUsersUserID(c *gin.Context) {
	sUserID := c.Param("userID")
	serviceUsers := service.NewUserService()

	userID, err := manager.SToUint64(sUserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx := server.Db.Begin()

	if err := serviceUsers.Delete(userID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	c.JSON(204, nil)
}

// private -------------------------------------------------------------------------------------------------------------
