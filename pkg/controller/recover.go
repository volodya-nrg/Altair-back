package controller

import (
	"altair/api/request"
	"altair/pkg/emailer"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"unicode/utf8"
)

// PostRecoverSendHash - создание записи/хеша для сверки, после, е-мэйла
func PostRecoverSendHash(c *gin.Context) {
	postRequest := new(request.PostRecoverySendHash)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	email := strings.TrimSpace(postRequest.Email)

	if !manager.ValidateEmail(email) {
		c.JSON(400, manager.ErrEmailNotCorrect.Error())
		return
	}

	serviceUsers := service.NewUserService()
	serviceRecovery := service.NewRecoveryService()

	// если пытаются восстановить пароль для соц. сети (емэйла), то выдадим ошибку
	if manager.IsSocialEmail(email) {
		c.JSON(400, manager.ErrEmailNotCorrect.Error())
		return
	}

	user, err := serviceUsers.GetUserByEmail(email)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return
	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	} else if !user.IsEmailConfirmed {
		c.JSON(400, manager.ErrEmailNotConfirmed.Error())
		return
	}

	hash := manager.RandASCII(manager.HashLen)
	link := fmt.Sprintf("%s/%s/%s", manager.Domain, "recover/check", hash)

	recovery := new(storage.Recovery)
	recovery.UserID = user.UserID
	recovery.Hash = hash

	if err := serviceRecovery.Create(recovery, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	emailData := emailer.Recover{Link: link}
	emailRequest := emailer.NewEmailRequest(email, "Восстановление пароля", "")
	if err := emailRequest.ParseTemplate("recover.html", emailData); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}
	if ok, err := emailRequest.SendMail(); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	} else if !ok {
		c.JSON(500, manager.ErrEmailNotSend.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204
}

// PostRecoverChangePass - создание обновленного пароля
func PostRecoverChangePass(c *gin.Context) {
	postRequest := new(request.PostRecoveryChangePass)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	hash := strings.TrimSpace(postRequest.Hash)
	password := strings.TrimSpace(postRequest.Password)
	passwordConfirm := strings.TrimSpace(postRequest.PasswordConfirm)

	if utf8.RuneCountInString(hash) < manager.HashLen {
		c.JSON(400, manager.ErrHashIsNotCorrect.Error())
		return
	}
	if utf8.RuneCountInString(password) < manager.MinLenPassword {
		c.JSON(400, manager.ErrPasswordIsShort.Error())
		return
	}
	if utf8.RuneCountInString(passwordConfirm) < manager.MinLenPassword {
		c.JSON(400, manager.ErrPasswordIsShort.Error())
		return
	}
	if password != passwordConfirm {
		c.JSON(400, manager.ErrPasswordsAreNotEqual.Error())
		return
	}

	serviceUser := service.NewUserService()
	serviceRecover := service.NewRecoveryService()

	recovery, err := serviceRecover.GetByHash(hash)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(400, err.Error())
		return
	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	user, err := serviceUser.GetUserByID(recovery.UserID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(400, err.Error())
		return
	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	user.Password = manager.HashAndSalt(password)

	if err := serviceUser.Update(user, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	if err := serviceRecover.Delete(user.UserID, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204
}

// private -------------------------------------------------------------------------------------------------------------
