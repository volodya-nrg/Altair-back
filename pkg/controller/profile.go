package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/configs"
	"altair/pkg/emailer"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/mediafire"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// PostProfile - создание профиля
func PostProfile(c *gin.Context) {
	postRequest := new(request.PostProfile)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceUsers := service.NewUserService()

	password := strings.TrimSpace(postRequest.Password)
	passwordConfirm := strings.TrimSpace(postRequest.PasswordConfirm)

	if password != passwordConfirm {
		c.JSON(400, manager.ErrPasswordsAreNotEqual.Error())
		return
	}

	if utf8.RuneCountInString(password) < manager.MinLenPassword {
		c.JSON(400, manager.ErrPasswordIsShort.Error())
		return
	}

	email := postRequest.Email
	if !manager.ValidateEmail(email) || manager.IsSocialEmail(email) {
		c.JSON(400, manager.ErrEmailNotCorrect.Error())
		return
	}

	user, err := serviceUsers.GetUserByEmail(email)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}
	if err == nil && user.IsEmailConfirmed {
		c.JSON(400, manager.ErrInterInYourProfile.Error())
		return
	}
	if gorm.IsRecordNotFoundError(err) {
		userNew := new(storage.User)
		userNew.Email = email
		userNew.Password = manager.HashAndSalt(password)

		if err := serviceUsers.Create(userNew, nil); err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}

		user = userNew
	}

	hash := manager.RandASCII(manager.HashLen)
	link := fmt.Sprintf("%s/%s/%s", manager.Domain, "check-email-through", hash) // ссылка на фронт

	user.HashForCheckEmail = hash
	if err := serviceUsers.Update(user, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	emailData := emailer.CheckEmail{Link: link}
	emailRequest := emailer.NewEmailRequest(email, "Верификация е-мэйла", "")
	if err := emailRequest.ParseTemplate("check_email.html", emailData); err != nil {
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

	c.JSON(http.StatusCreated, user) // 201 - создана новая запись
}

// PutProfile - редактирование профиля
func PutProfile(c *gin.Context) {
	putRequest := new(request.PutProfile)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceUsers := service.NewUserService()
	serviceMediafire := mediafire.NewMediafireService()
	servicePhone := service.NewPhoneService()

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
	passwordOld := strings.TrimSpace(putRequest.PasswordOld)
	passwordNew := strings.TrimSpace(putRequest.PasswordNew)
	passwordConfirm := strings.TrimSpace(putRequest.PasswordConfirm)

	if utf8.RuneCountInString(passwordOld) > 0 && utf8.RuneCountInString(passwordNew) > 0 && utf8.RuneCountInString(passwordConfirm) > 0 {
		if passwordNew != passwordConfirm {
			c.JSON(400, manager.ErrPasswordsAreNotEqual.Error())
			return
		}
		if !manager.ComparePasswords(user.Password, passwordOld) {
			c.JSON(400, manager.ErrPasswordsAreNotEqualCurrentAndOld.Error())
			return
		}

		user.Password = manager.HashAndSalt(passwordNew)
	}

	var filePath string
	if len(form.File["files"]) > 0 {
		file := form.File["files"][0] // только одно фото

		fileName, err := manager.UploadImage(file, manager.DirImages, c.SaveUploadedFile)
		if err != nil {
			logger.Warning.Println(err.Error())

		} else if fileName != "" {
			if externalFilePath, err := serviceMediafire.UploadSimple(manager.DirImages + "/" + fileName); err != nil {
				logger.Warning.Println(err.Error())

			} else if externalFilePath == "" {
				logger.Warning.Println(manager.ErrNotFoundExternalFilePath.Error())

			} else {
				if err := os.Remove(manager.DirImages + "/" + fileName); err != nil {
					logger.Warning.Println(err.Error())
				}

				filePath = externalFilePath
			}
		}
	}

	// отправили файл или у Юзера есть аватар и пришедший аватар пустой, удалим
	if filePath != "" || (user.Avatar != "" && putRequest.Avatar == "") {
		if user.Avatar != "" {
			// тут по идеи надо удалить файл с внешнего хранилища
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

	phones, err := servicePhone.GetPhonesByUserID(userID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	userExt := new(response.UserExt)
	userExt.User = user
	userExt.Phones = phones

	c.JSON(200, userExt)
}

// DeleteProfile - удаление профиля
func DeleteProfile(c *gin.Context) {
	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceAds := service.NewAdService()
	serviceUser := service.NewUserService()
	servicePhone := service.NewPhoneService()

	tx := server.Db.Begin()

	// удаляем профиль
	if err := serviceUser.Delete(userID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// удаляем его объявления
	if err := serviceAds.DeleteAllByUserID(userID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// удаляем его номера телефонов
	if err := servicePhone.DeleteAllByUserID(userID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	c.JSON(204, nil)
}

// GetProfileAds - получение всех объявлений профиля
func GetProfileAds(c *gin.Context) {
	req := new(request.GetProfileAds)

	if err := c.ShouldBind(req); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceAds := service.NewAdService()

	ads, err := serviceAds.GetAdsFullByUserID(userID, "created_at", req.Limit, req.Offset)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, ads)
}

// GetProfileAdsAdID - получение конкретного объявления профиля
func GetProfileAdsAdID(c *gin.Context) {
	sAdID := c.Param("adID")
	adID, err := manager.SToUint64(sAdID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceAds := service.NewAdService()

	adFull, err := serviceAds.GetAdFullByUserIDAndCatID(userID, adID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return
	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, adFull)
}

// GetProfileSettings - получение настроение профиля
func GetProfileSettings(c *gin.Context) {
	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
	}
	//----------------------------------------------------------

	c.JSON(200, userID)
}

// GetProfileInfo - получение основной информации профиля
func GetProfileInfo(c *gin.Context) {
	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	serviceUsers := service.NewUserService()
	servicePhone := service.NewPhoneService()

	user, err := serviceUsers.GetUserByID(userID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	phones, err := servicePhone.GetPhonesByUserID(userID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	userExt := new(response.UserExt)
	userExt.User = user
	userExt.Phones = phones

	c.JSON(200, userExt)
}

// GetProfileCheckEmailThroughHash - верификация е-мэйла. Проверка е-мэйла через хеш.
func GetProfileCheckEmailThroughHash(c *gin.Context) {
	hash := c.Param("hash")
	serviceUsers := service.NewUserService()

	user, err := serviceUsers.GetUserByHashCheckEmail(hash)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(400, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	user.IsEmailConfirmed = true
	user.HashForCheckEmail = ""

	if err := serviceUsers.Update(user, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(204, nil)
}

// PostProfilePhone - добавление номера телефона для профиля
func PostProfilePhone(c *gin.Context) {
	postRequest := new(request.PostProfilePhone)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	servicePhone := service.NewPhoneService()
	serviceSMS := service.NewSMSService(configs.Cfg.SMS.APIKey, configs.Cfg.SMS.Domain)

	re := regexp.MustCompile(`\d+`)
	slicePhone := re.FindAllString(postRequest.Number, -1)
	number := strings.Join(slicePhone, "")

	if matched, err := regexp.Match(manager.PhonePattern, []byte(number)); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return

	} else if !matched {
		c.JSON(400, manager.ErrPhoneNotCorrect.Error())
		return
	}

	phoneOld, err := servicePhone.GetPhoneByNumberAndUserID(number, userID)
	if !gorm.IsRecordNotFoundError(err) {
		if phoneOld.IsVerify {
			c.JSON(400, manager.ErrPhoneAlreadyVerified.Error())
			return
		}

		deltaSec := time.Since(phoneOld.CreatedAt).Seconds()
		if deltaSec < manager.MinSecBetweenSendSmsForVerifyPhone {
			c.JSON(400, manager.ErrPhoneTimeLimit.Error())
			return
		}
		if err := servicePhone.Delete(number, userID, true, nil); err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}
	}

	phone := new(storage.Phone)
	phone.Number = number
	phone.UserID = userID
	phone.Code = fmt.Sprintf("%d", manager.RandIntByRange(10000, 99999))
	messageForSMS := "От: www.Altair.uz\nВаш код подтверждения: " + phone.Code + ".\nНаберите его в поле ввода."

	respBalance, err := serviceSMS.Balance()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	requestCost := service.SMSCostRequest{
		To:  phone.Number,
		Msg: messageForSMS,
	}
	if resp, err := serviceSMS.Cost(requestCost); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	} else if resp.Status == "OK" {
		// 2.79 руб. для России, 4.18 руб. для Узб
		for _, v := range resp.Sms {
			if v.Status == "ERROR" {
				if v.StatusText != "" {
					logger.Warning.Println(v.StatusText)
				}

				c.JSON(400, errors.New(serviceSMS.GetStatusInfo(v.StatusCode)).Error())
				return
			}
			if v.Status == "OK" && respBalance.Balance < v.Cost {
				logger.Warning.Println(manager.ErrSMSMoneyInBalanceIsSmall.Error())
				c.JSON(500, manager.ErrSMSSendingIsNotYetPossible.Error())
				return
			}
		}
	}

	requestToSMS := service.SMSSendRequest{
		To:  phone.Number,
		Msg: messageForSMS,
	}
	if resp, err := serviceSMS.Send(requestToSMS); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return

	} else if resp.Status == "OK" {
		for _, v := range resp.Sms {
			if v.Status == "ERROR" {
				if v.StatusText != "" {
					logger.Warning.Println(v.StatusText)
				}

				c.JSON(400, errors.New(serviceSMS.GetStatusInfo(v.StatusCode)).Error())
				return
			}
		}
	}

	if err := servicePhone.Create(phone, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusCreated, phone) // 201 - создана новая запись
}

// PutProfilePhoneNumberCode - проверка номера тел. профиля через проверочный код
func PutProfilePhoneNumberCode(c *gin.Context) {
	number := c.Param("number")
	code := c.Param("code")

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	servicePhone := service.NewPhoneService()
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

	phone, err := servicePhone.GetPhoneByNumberAndUserID(number, user.UserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// Если прошло определенное время, то проверочный код уже не актуален
	deltaSec := time.Since(phone.CreatedAt).Seconds()
	if deltaSec > manager.MinSecLifeVerifyCode {
		c.JSON(400, manager.ErrPhoneCodeIsNotActual.Error())
		return
	}

	if phone.Code != code {
		c.JSON(400, manager.ErrPhoneNotCorrectCode.Error())
		return
	}

	phone.IsVerify = true

	if err := servicePhone.Update(phone, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	phones, err := servicePhone.GetPhonesByUserID(user.UserID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// отдаем профиля, чтоб везде "встали" корректные номера телефонов
	userExt := new(response.UserExt)
	userExt.User = user
	userExt.Phones = phones

	c.JSON(http.StatusOK, userExt)
}

// DeleteProfilePhoneNumber - удаление номера телефона
func DeleteProfilePhoneNumber(c *gin.Context) {
	number := c.Param("number")

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(http.StatusForbidden, manager.ErrUndefinedUserID.Error()) // 403
		return
	}
	//----------------------------------------------------------

	servicePhone := service.NewPhoneService()
	serviceAds := service.NewAdService()

	var phoneIDNew uint64 = 0
	curPhone := new(storage.Phone)
	phones, err := servicePhone.GetPhonesByUserID(userID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// подхватим нужные данные
	for _, v := range phones {
		if v.Number == number {
			curPhone = v
		} else {
			phoneIDNew = v.PhoneID
		}
	}

	tx := server.Db.Begin()

	if err := servicePhone.Delete(curPhone.Number, curPhone.UserID, false, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// если номер один, то все объявления, приуроченные к этому номеру, установить статус = отключено, phoneId=0
	// если неск-ко валидных номеров тел., то возьмем валидный
	if err := serviceAds.UpdateByPhoneID(curPhone.PhoneID, phoneIDNew, userID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	c.JSON(204, nil)
}
