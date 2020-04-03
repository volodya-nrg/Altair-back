package controller

import (
	"altair/api/request"
	"altair/api/response"
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
	pPostRequest := new(request.PostUser)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := postUsers(pPostRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
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

	res := putUsersUserId(c.Param("userId"), pPutRequest, form, c.SaveUploadedFile)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
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
func postUsers(pPostRequest *request.PostUser) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	if pPostRequest.Password != pPostRequest.PasswordConfirm {
		res.Status = 400
		res.Err = errPasswordsAreNotEqual
		return res
	}

	user := new(storage.User)
	user.Email = pPostRequest.Email
	user.Password = pPostRequest.Password

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
func putUsersUserId(sUserId string, pPutRequest *request.PutUser, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) response.Result {
	serviceUsers := service.NewUserService()
	res := response.Result{}

	userId, err := strconv.ParseUint(sUserId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pUser, err := serviceUsers.GetUserByID(userId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 500
		res.Err = err
		return res
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
			res.Status = 400
			res.Err = errPasswordsAreNotEqual
			return res
		}
		if !helpers.ComparePasswords(pUser.Password, passwordOld) {
			res.Status = 400
			res.Err = errPasswordsAreNotEqualCurrentAndOld
			return res
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

	if err = serviceUsers.Update(pUser, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pUser
	return res
}
