package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	errPasswordsAreNotEqual              = errors.New("passwords are not equal")
	errPasswordsAreNotEqualCurrentAndOld = errors.New("passwords are not equal (current, old)")
)

func GetUsers(c *gin.Context) {
	res := getUsers()
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetUsersUserId(c *gin.Context) {
	res := getUsersUserId(c.Param("userId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostUsers(c *gin.Context) {
	postRequest := new(request.PostUser)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := postUsers(postRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutUsersUserId(c *gin.Context) {
	putRequest := new(request.PutUser)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err)
		c.JSON(500, err.Error())
		return
	}

	res := putUsersUserId(c.Param("userId"), putRequest, form, c.SaveUploadedFile)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeleteUsersUserId(c *gin.Context) {
	res := deleteUsersUserId(c.Param("userId"))
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getUsers() response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	users, err := serviceUsers.GetUsers()
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = users
	return res
}
func getUsersUserId(sUserId string) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	user, err := serviceUsers.GetUserByID(userId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = user
	return res
}
func postUsers(postRequest *request.PostUser) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	if postRequest.Password != postRequest.PasswordConfirm {
		res.Status = 400
		res.Err = errPasswordsAreNotEqual
		return res
	}

	user := new(storage.User)
	user.Email = postRequest.Email
	user.Password = postRequest.Password

	if err := serviceUsers.Create(user, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = user
	return res
}
func putUsersUserId(sUserId string, putRequest *request.PutUser, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	user, err := serviceUsers.GetUserByID(userId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	passwordOld := strings.TrimSpace(putRequest.PasswordOld)
	password := strings.TrimSpace(putRequest.Password)
	passwordConfirm := strings.TrimSpace(putRequest.PasswordConfirm)

	user.Name = strings.TrimSpace(putRequest.Name)

	if user.IsEmailConfirmed != putRequest.IsEmailConfirmed {
		user.IsEmailConfirmed = !user.IsEmailConfirmed
	}

	if utf8.RuneCountInString(passwordOld) > 0 && utf8.RuneCountInString(password) > 0 && utf8.RuneCountInString(passwordConfirm) > 0 {
		if password != passwordConfirm {
			res.Status = 400
			res.Err = errPasswordsAreNotEqual
			return res
		}
		if !helpers.ComparePasswords(user.Password, passwordOld) {
			res.Status = 400
			res.Err = errPasswordsAreNotEqualCurrentAndOld
			return res
		}
		user.Password = helpers.HashAndSalt(password)
	}

	var filePath string
	for _, file := range form.File["file"] {
		filePath, err = helpers.UploadImage(file, "./web/images", fnUpload)
		if err != nil {
			logger.Warning.Println(err)
		}
		break // только один файл
	}

	// отправили файл или у Юзера есть аватар и пришедший аватар пустой (удалим)
	if filePath != "" || (user.Avatar != "" && putRequest.Avatar == "") {
		if user.Avatar != "" {
			tmp := "./web/images/" + user.Avatar

			if helpers.FileExists(tmp) {
				if err := os.Remove(tmp); err != nil {
					logger.Warning.Println(err)
				}
			}

			user.Avatar = ""
		}
		if filePath != "" {
			user.Avatar = filePath
		}
	}

	if err = serviceUsers.Update(user, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = user
	return res
}
func deleteUsersUserId(sUserId string) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()
	if err := serviceUsers.Delete(userId, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}
	tx.Commit()

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
