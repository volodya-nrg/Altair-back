package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

var (
	errAlreadyHasReservedName = errors.New("already name exists")
)

func GetProperties(c *gin.Context) {
	res := getProperties(c.DefaultQuery("catId", ""))
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
	postRequest := new(request.PostProperty)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	// тут нужно проверить чтоб не было зарезирвированных уже полей

	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := postProperties(postRequest, mValueId, mValueTitle, mValuePos)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutPropertiesPropertyId(c *gin.Context) {
	putRequest := new(request.PutProperty)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	propertyId := c.Param("propertyId")
	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := putPropertiesPropertyId(propertyId, putRequest, mValueId, mValueTitle, mValuePos)
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
func getProperties(catIdSrc string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperty := service.NewValuesPropertyService()
	res := response.Result{}

	if catIdSrc == "" {
		catIdSrc = "0"
	}

	catId, err := strconv.ParseUint(catIdSrc, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}
	if catId > 0 {
		propertiesFull, err := serviceProperties.GetPropertiesFullByCatId(catId, false, serviceValuesProperty)
		if err != nil {
			logger.Warning.Println(err)
			res.Status = 500
			res.Err = err
			return res
		}

		res.Status = 200
		res.Err = nil
		res.Data = propertiesFull
		return res
	}

	propertiesWithKindName, err := serviceProperties.GetPropertiesWithKindName()
	if err != nil {
		logger.Warning.Println(err)
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
func postProperties(postRequest *request.PostProperty,
	mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}
	prop := new(storage.Property)
	sliceJson := helpers.GetTagsFromStruct(storage.Ad{}, "json")

	if has, _ := helpers.InArray(strings.TrimSpace(postRequest.Name), sliceJson); has {
		res.Status = 400
		res.Err = errAlreadyHasReservedName
		return res
	}

	tx := server.Db.Debug().Begin()

	prop.Title = strings.TrimSpace(postRequest.Title)
	prop.KindPropertyId = postRequest.KindPropertyId
	prop.Name = strings.TrimSpace(postRequest.Name)
	prop.Suffix = strings.TrimSpace(postRequest.Suffix)
	prop.Comment = strings.TrimSpace(postRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(postRequest.PrivateComment)

	if err := serviceProperties.Create(prop, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	// если данные приходят пачкой, через textarea, то перезапишем их в соот-ие map-ы
	selectAsTextarea := strings.TrimSpace(postRequest.SelectAsTextarea)
	if selectAsTextarea != "" {
		mId, mTitle, mPos = createMapsFromMultiText(selectAsTextarea)
	}

	_, err := serviceProperties.ReWriteValuesForProperties(prop.PropertyId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propertyFull, err := serviceProperties.GetPropertyFullById(prop.PropertyId, serviceValuesProperties)
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
func putPropertiesPropertyId(sPropertyId string, putRequest *request.PutProperty,
	mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}
	sliceJson := helpers.GetTagsFromStruct(storage.Ad{}, "json")

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	if has, _ := helpers.InArray(strings.TrimSpace(putRequest.Name), sliceJson); has {
		res.Status = 400
		res.Err = errAlreadyHasReservedName
		return res
	}

	prop, err := serviceProperties.GetPropertyById(propertyId)
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

	prop.Title = strings.TrimSpace(putRequest.Title)
	prop.Name = strings.TrimSpace(putRequest.Name)
	prop.KindPropertyId = putRequest.KindPropertyId
	prop.Suffix = strings.TrimSpace(putRequest.Suffix)
	prop.Comment = strings.TrimSpace(putRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(putRequest.PrivateComment)

	if err = serviceProperties.Update(prop, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	_, err = serviceProperties.ReWriteValuesForProperties(prop.PropertyId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propertyFull, err := serviceProperties.GetPropertyFullById(prop.PropertyId, serviceValuesProperties)
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
	serviceValueProperties := service.NewValuesPropertyService()
	res := response.Result{}

	propertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()
	if err := serviceProperties.Delete(propertyId, serviceValueProperties, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}
	tx.Commit()

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
func createMapsFromMultiText(multiText string) (map[string]string, map[string]string, map[string]string) {
	aStrings := strings.Split(multiText, "\n")
	mId := make(map[string]string, 0)
	mTitle := make(map[string]string, 0)
	mPos := make(map[string]string, 0)

	for k, v := range aStrings {
		a := fmt.Sprint(k + 1)

		mId[a] = "0"
		mTitle[a] = strings.TrimSpace(v)
		mPos[a] = a
	}

	return mId, mTitle, mPos
}
