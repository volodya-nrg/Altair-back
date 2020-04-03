package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetProperties(c *gin.Context) {
	res := getProperties()
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetPropertiesPropertyId(c *gin.Context) {
	res := getPropertiesPropertyId(c.Param("propertyId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostProperties(c *gin.Context) {
	pPostRequest := new(request.PostProperty)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := postProperties(pPostRequest, mValueId, mValueTitle, mValuePos)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutPropertiesPropertyId(c *gin.Context) {
	pPutRequest := new(request.PutProperty)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	propertyId := c.Param("propertyId")
	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := putPropertiesPropertyId(propertyId, pPutRequest, mValueId, mValueTitle, mValuePos)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeletePropertiesPropertyId(c *gin.Context) {
	res := deletePropertiesPropertyId(c.Param("propertyId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getProperties() response.Result {
	serviceProperties := service.NewPropertyService()
	res := response.Result{}

	propertiesWithKindName, err := serviceProperties.GetPropertiesWithKindName()
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = propertiesWithKindName
	return res
}
func getPropertiesPropertyId(sPropertyId string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	propertyFull, err := serviceProperties.GetPropertyFullById(propertyId, serviceValuesProperties)
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
	res.Data = propertyFull
	return res
}
func postProperties(pPostRequest *request.PostProperty, mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}
	pProperty := new(storage.Property)
	tx := server.Db.Debug().Begin()

	pProperty.Title = strings.TrimSpace(pPostRequest.Title)
	pProperty.KindPropertyId = pPostRequest.KindPropertyId
	pProperty.Name = strings.TrimSpace(pPostRequest.Name)

	if err := serviceProperties.Create(pProperty, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	_, err := serviceProperties.ReWriteValuesForProperties(pProperty.PropertyId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propertyFull, err := serviceProperties.GetPropertyFullById(pProperty.PropertyId, serviceValuesProperties)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = propertyFull
	return res
}
func putPropertiesPropertyId(sPropertyId string, putRequest *request.PutProperty, mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pProperty, err := serviceProperties.GetPropertyById(propertyId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()

	pProperty.Title = strings.TrimSpace(putRequest.Title)
	pProperty.Name = strings.TrimSpace(putRequest.Name)
	pProperty.KindPropertyId = putRequest.KindPropertyId

	if err = serviceProperties.Update(pProperty, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	_, err = serviceProperties.ReWriteValuesForProperties(pProperty.PropertyId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propertyFull, err := serviceProperties.GetPropertyFullById(pProperty.PropertyId, serviceValuesProperties)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = propertyFull
	return res
}
func deletePropertiesPropertyId(sPropertyId string) response.Result {
	serviceProperties := service.NewPropertyService()
	res := response.Result{}

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceProperties.Delete(propertyId, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
