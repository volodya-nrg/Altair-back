package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

// GetKindProps - получить все разновидности свойств
func GetKindProps(c *gin.Context) {
	serviceKindProps := service.NewKindPropService()

	pKindProps, err := serviceKindProps.GetKindProps("kind_prop_id desc")
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, pKindProps)
}

// GetKindPropsKindPropID - получить разновидности конкретного свойства
func GetKindPropsKindPropID(c *gin.Context) {
	sKindPropID := c.Param("kindPropID")
	serviceKindProps := service.NewKindPropService()

	kindPropID, err := manager.SToUint64(sKindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	pKindProp, err := serviceKindProps.GetKindPropByID(kindPropID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, pKindProp)
}

// PostKindProps - создание вида св-ва
func PostKindProps(c *gin.Context) {
	pPostRequest := new(request.PostKindProp)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceKindProps := service.NewKindPropService()
	pKindProp := new(storage.KindProp)

	pKindProp.Name = pPostRequest.Name

	if err := serviceKindProps.Create(pKindProp, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(201, pKindProp)
}

// PutKindPropsKindPropID - изменение одного конкретного вида свойства
func PutKindPropsKindPropID(c *gin.Context) {
	sKindPropID := c.Param("kindPropID")
	putRequest := new(request.PutKindProp)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceKindProps := service.NewKindPropService()

	kindPropID, err := manager.SToUint64(sKindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	pKindProp, err := serviceKindProps.GetKindPropByID(kindPropID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	pKindProp.Name = strings.TrimSpace(putRequest.Name)

	if err = serviceKindProps.Update(pKindProp, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, pKindProp)
}

// DeleteKindPropsKindPropID - удаление вида свойства
func DeleteKindPropsKindPropID(c *gin.Context) {
	sKindPropID := c.Param("kindPropID")
	serviceKindProps := service.NewKindPropService()

	kindPropID, err := manager.SToUint64(sKindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	if err := serviceKindProps.Delete(kindPropID, nil); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(204, nil)
}

// private -------------------------------------------------------------------------------------------------------------
