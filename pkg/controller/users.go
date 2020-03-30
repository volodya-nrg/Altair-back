package controller

import (
	"altair/api/request"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
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
	pResult := getUsers()
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func GetUsersUserId(c *gin.Context) {
	pResult := getUsersUserId(c.Param("userId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PostUsers(c *gin.Context) {
	pPostRequest := new(request.PostUser)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := postUsers(pPostRequest)
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PutUsersUserId(c *gin.Context) {
	pPutRequest := new(request.PutUser)

	if err := c.ShouldBind(pPutRequest); err != nil {
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

	pResult := putUsersUserId(c.Param("userId"), pPutRequest, form, c.SaveUploadedFile)
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getUsers() *result {
	serviceUsers := service.NewUserService()
	pResult := new(result)
	users, err := serviceUsers.GetUsers()

	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = users
	return pResult
}
func getUsersUserId(sUserId string) *result {
	serviceUsers := service.NewUserService()
	pResult := new(result)

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	user, err := serviceUsers.GetUserByID(userId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = user
	return pResult
}
func postUsers(pPostRequest *request.PostUser) *result {
	serviceUsers := service.NewUserService()
	pResult := new(result)

	if pPostRequest.Password != pPostRequest.PasswordConfirm {
		pResult.Status = 400
		pResult.Err = errPasswordsAreNotEqual
		return pResult
	}

	user := new(storage.User)
	user.Email = pPostRequest.Email
	user.Password = pPostRequest.Password

	if err := serviceUsers.Create(user); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 201
	pResult.Err = nil
	pResult.Data = user
	return pResult
}
func putUsersUserId(sUserId string, pPutRequest *request.PutUser, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) *result {
	serviceUsers := service.NewUserService()
	pResult := new(result)

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pUser, err := serviceUsers.GetUserByID(userId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	passwordOld := strings.TrimSpace(pPutRequest.PasswordOld)
	password := strings.TrimSpace(pPutRequest.Password)
	passwordConfirm := strings.TrimSpace(pPutRequest.PasswordConfirm)

	pUser.Name = strings.TrimSpace(pPutRequest.Name)

	if pUser.EmailIsConfirmed != pPutRequest.EmailIsConfirmed {
		pUser.EmailIsConfirmed = !pUser.EmailIsConfirmed
	}

	if utf8.RuneCountInString(passwordOld) > 0 && utf8.RuneCountInString(password) > 0 && utf8.RuneCountInString(passwordConfirm) > 0 {
		if password != passwordConfirm {
			pResult.Status = 400
			pResult.Err = errPasswordsAreNotEqual
			return pResult
		}
		if !helpers.ComparePasswords(pUser.Password, passwordOld) {
			pResult.Status = 400
			pResult.Err = errPasswordsAreNotEqualCurrentAndOld
			return pResult
		}
		pUser.Password = helpers.HashAndSalt(password)
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
	if filePath != "" || (pUser.Avatar != "" && pPutRequest.Avatar == "") {
		if pUser.Avatar != "" {
			tmp := "./web/images/" + pUser.Avatar

			if helpers.FileExists(tmp) {
				if err := os.Remove(tmp); err != nil {
					logger.Warning.Println(err)
				}
			}

			pUser.Avatar = ""
		}
		if filePath != "" {
			pUser.Avatar = filePath
		}
	}

	if err = serviceUsers.Update(pUser); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = pUser
	return pResult
}
