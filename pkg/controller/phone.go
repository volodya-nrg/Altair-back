package controller

import (
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetPhonesPhoneID - получение данных об конкретном номере телефона
func GetPhonesPhoneID(c *gin.Context) {
	sPhoneID := c.Param("phoneID")
	servicePhone := service.NewPhoneService()

	phoneID, err := manager.SToUint64(sPhoneID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	phone, err := servicePhone.GetByID(phoneID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, phone)
}

// private -------------------------------------------------------------------------------------------------------------
