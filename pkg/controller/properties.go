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

func GetProperties(c *gin.Context) {
	pResult := getProperties()
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func GetPropertiesPropertyId(c *gin.Context) {
	pResult := getPropertiesPropertyId(c.Param("propertyId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PostProperties(c *gin.Context) {
	pPostRequest := new(request.PostProperty)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := postProperties(pPostRequest, c.PostFormMap("valueId"), c.PostFormMap("valueTitle"), c.PostFormMap("valuePos"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PutPropertiesPropertyId(c *gin.Context) {
	pPutRequest := new(request.PutProperty)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := putPropertiesPropertyId(
		c.Param("propertyId"),
		pPutRequest,
		c.PostFormMap("valueId"),
		c.PostFormMap("valueTitle"),
		c.PostFormMap("valuePos"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func DeletePropertiesPropertyId(c *gin.Context) {
	pResult := deletePropertiesPropertyId(c.Param("propertyId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getProperties() *result {
	serviceProperties := service.NewPropertyService()
	pResult := new(result)

	propertiesFull, err := serviceProperties.GetPropertiesFull()
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = propertiesFull
	return pResult
}
func getPropertiesPropertyId(sPropertyId string) *result {
	serviceProperties := service.NewPropertyService()
	pResult := new(result)

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	propertyFull, err := serviceProperties.GetPropertyFullById(propertyId)
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
	pResult.Data = propertyFull
	return pResult
}
func postProperties(pPostRequest *request.PostProperty, mId map[string]string, mTitle map[string]string, mPos map[string]string) *result {
	serviceProperties := service.NewPropertyService()
	pResult := new(result)
	pProperty := new(storage.Property)

	pProperty.Title = strings.TrimSpace(pPostRequest.Title)
	pProperty.KindPropertyId = pPostRequest.KindPropertyId
	pProperty.Name = strings.TrimSpace(pPostRequest.Name)
	pProperty.IsCanAsFilter = pPostRequest.IsCanAsFilter
	pProperty.MaxInt = pPostRequest.MaxInt

	if err := serviceProperties.Create(pProperty); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	_, err := serviceProperties.ReWriteValuesForProperties(pProperty.PropertyId, mId, mTitle, mPos)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	propertyFull, err := serviceProperties.GetPropertyFullById(pProperty.PropertyId)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 201
	pResult.Err = nil
	pResult.Data = propertyFull
	return pResult
}
func putPropertiesPropertyId(sPropertyId string, putRequest *request.PutProperty, mId map[string]string, mTitle map[string]string, mPos map[string]string) *result {
	serviceProperties := service.NewPropertyService()
	pResult := new(result)

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pProperty, err := serviceProperties.GetPropertyById(propertyId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pProperty.Title = strings.TrimSpace(putRequest.Title)
	pProperty.Name = strings.TrimSpace(putRequest.Name)
	pProperty.KindPropertyId = putRequest.KindPropertyId
	pProperty.MaxInt = putRequest.MaxInt
	pProperty.IsCanAsFilter = putRequest.IsCanAsFilter

	if err = serviceProperties.Update(pProperty); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	_, err = serviceProperties.ReWriteValuesForProperties(pProperty.PropertyId, mId, mTitle, mPos)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	propertyFull, err := serviceProperties.GetPropertyFullById(pProperty.PropertyId)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = propertyFull
	return pResult
}
func deletePropertiesPropertyId(sPropertyId string) *result {
	serviceProperties := service.NewPropertyService()
	pResult := new(result)

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	if err := serviceProperties.Delete(propertyId); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 204
	pResult.Err = nil
	pResult.Data = nil
	return pResult
}
