package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetKindProperties(c *gin.Context) {
	pResult := getKindProperties()
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func GetKindPropertiesKindPropertyId(c *gin.Context) {
	pResult := getKindPropertiesKindPropertyId(c.Param("kindPropertyId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PostKindProperties(c *gin.Context) {
	pPostRequest := new(request.PostKindProperty)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := postKindProperties(pPostRequest)
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PutKindPropertiesKindPropertyId(c *gin.Context) {
	pPutRequest := new(request.PutKindProperty)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := putKindPropertiesKindPropertyId(c.Param("kindPropertyId"), pPutRequest)
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func DeleteKindPropertiesKindPropertyId(c *gin.Context) {
	pResult := deleteKindPropertiesKindPropertyId(c.Param("kindPropertyId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getKindProperties() *result {
	serviceKindProperties := service.NewKindPropertyService()
	pResult := new(result)

	pKindProperties, err := serviceKindProperties.GetKindProperties()
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = pKindProperties
	return pResult
}
func getKindPropertiesKindPropertyId(sKindPropertyId string) *result {
	serviceKindProperties := service.NewKindPropertyService()
	pResult := new(result)

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pKindProperty, err := serviceKindProperties.GetKindPropertyById(kindPropertyId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = pKindProperty
	return pResult
}
func postKindProperties(pPostRequest *request.PostKindProperty) *result {
	serviceKindProperties := service.NewKindPropertyService()
	pResult := new(result)
	pKindProperty := new(storage.KindProperty)

	pKindProperty.Name = pPostRequest.Name

	if err := serviceKindProperties.Create(pKindProperty); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pResult.Status = 201
	pResult.Err = nil
	pResult.Data = pKindProperty
	return pResult
}
func putKindPropertiesKindPropertyId(sKindPropertyId string, putRequest *request.PutKindProperty) *result {
	serviceKindProperties := service.NewKindPropertyService()
	pResult := new(result)

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pKindProperty, err := serviceKindProperties.GetKindPropertyById(kindPropertyId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pKindProperty.Name = strings.TrimSpace(putRequest.Name)

	if err = serviceKindProperties.Update(pKindProperty); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = pKindProperty
	return pResult
}
func deleteKindPropertiesKindPropertyId(sKindPropertyId string) *result {
	serviceKindProperties := service.NewKindPropertyService()
	pResult := new(result)

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	if err := serviceKindProperties.Delete(kindPropertyId); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 204
	pResult.Err = nil
	pResult.Data = nil
	return pResult
}
