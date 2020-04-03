package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetKindProperties(c *gin.Context) {
	res := getKindProperties()
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetKindPropertiesKindPropertyId(c *gin.Context) {
	res := getKindPropertiesKindPropertyId(c.Param("kindPropertyId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostKindProperties(c *gin.Context) {
	pPostRequest := new(request.PostKindProperty)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := postKindProperties(pPostRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutKindPropertiesKindPropertyId(c *gin.Context) {
	pPutRequest := new(request.PutKindProperty)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := putKindPropertiesKindPropertyId(c.Param("kindPropertyId"), pPutRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeleteKindPropertiesKindPropertyId(c *gin.Context) {
	res := deleteKindPropertiesKindPropertyId(c.Param("kindPropertyId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getKindProperties() response.Result {
	serviceKindProperties := service.NewKindPropertyService()
	res := response.Result{}

	pKindProperties, err := serviceKindProperties.GetKindProperties()
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pKindProperties
	return res
}
func getKindPropertiesKindPropertyId(sKindPropertyId string) response.Result {
	serviceKindProperties := service.NewKindPropertyService()
	res := response.Result{}

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	pKindProperty, err := serviceKindProperties.GetKindPropertyById(kindPropertyId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pKindProperty
	return res
}
func postKindProperties(pPostRequest *request.PostKindProperty) response.Result {
	serviceKindProperties := service.NewKindPropertyService()
	res := response.Result{}
	pKindProperty := new(storage.KindProperty)

	pKindProperty.Name = pPostRequest.Name

	if err := serviceKindProperties.Create(pKindProperty, nil); err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = pKindProperty
	return res
}
func putKindPropertiesKindPropertyId(sKindPropertyId string, putRequest *request.PutKindProperty) response.Result {
	serviceKindProperties := service.NewKindPropertyService()
	res := response.Result{}

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pKindProperty, err := serviceKindProperties.GetKindPropertyById(kindPropertyId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	pKindProperty.Name = strings.TrimSpace(putRequest.Name)

	if err = serviceKindProperties.Update(pKindProperty, nil); err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pKindProperty
	return res
}
func deleteKindPropertiesKindPropertyId(sKindPropertyId string) response.Result {
	serviceKindProperties := service.NewKindPropertyService()
	res := response.Result{}

	kindPropertyId, err := strconv.ParseUint(sKindPropertyId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceKindProperties.Delete(kindPropertyId, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
