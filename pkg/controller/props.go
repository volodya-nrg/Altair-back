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

func GetProps(c *gin.Context) {
	res := getProps(c.DefaultQuery("catId", ""))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetPropsPropId(c *gin.Context) {
	res := getPropsPropId(c.Param("propId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostProps(c *gin.Context) {
	postRequest := new(request.PostProp)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	// тут нужно проверить чтоб не было зарезирвированных уже полей

	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := postProps(postRequest, mValueId, mValueTitle, mValuePos)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutPropsPropId(c *gin.Context) {
	putRequest := new(request.PutProp)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	propId := c.Param("propId")
	mValueId := c.PostFormMap("valueId")
	mValueTitle := c.PostFormMap("valueTitle")
	mValuePos := c.PostFormMap("valuePos")

	res := putPropsPropId(propId, putRequest, mValueId, mValueTitle, mValuePos)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeletePropsPropId(c *gin.Context) {
	res := deletePropsPropId(c.Param("propId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getProps(catIdSrc string) response.Result {
	serviceProps := service.NewPropService()
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
	// вытаскивать propFull слишком накладно (1.4мб)
	if catId > 0 {
		propsFull, err := serviceProps.GetPropsFullByCatId(catId, false)
		if err != nil {
			logger.Warning.Println(err)
			res.Status = 500
			res.Err = err
			return res
		}

		res.Status = 200
		res.Err = nil
		res.Data = propsFull
		return res
	}

	propsWithKindName, err := serviceProps.GetPropsWithKindName()
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = propsWithKindName
	return res
}
func getPropsPropId(sPropId string) response.Result {
	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()
	res := response.Result{}

	propId, err := strconv.ParseUint(sPropId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	propFull, err := serviceProps.GetPropFullById(propId, serviceValuesProps)
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
	res.Data = propFull
	return res
}
func postProps(postRequest *request.PostProp,
	mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()
	res := response.Result{}
	prop := new(storage.Prop)
	sliceJson := helpers.GetTagsFromStruct(storage.Ad{}, "json")

	if has, _ := helpers.InArray(strings.TrimSpace(postRequest.Name), sliceJson); has {
		res.Status = 400
		res.Err = errAlreadyHasReservedName
		return res
	}

	tx := server.Db.Debug().Begin()

	prop.Title = strings.TrimSpace(postRequest.Title)
	prop.KindPropId = postRequest.KindPropId
	prop.Name = strings.TrimSpace(postRequest.Name)
	prop.Suffix = strings.TrimSpace(postRequest.Suffix)
	prop.Comment = strings.TrimSpace(postRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(postRequest.PrivateComment)

	if err := serviceProps.Create(prop, tx); err != nil {
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

	_, err := serviceProps.ReWriteValuesForProps(prop.PropId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propFull, err := serviceProps.GetPropFullById(prop.PropId, serviceValuesProps)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = propFull
	return res
}
func putPropsPropId(sPropId string, putRequest *request.PutProp,
	mId map[string]string, mTitle map[string]string, mPos map[string]string) response.Result {
	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()
	res := response.Result{}
	sliceJson := helpers.GetTagsFromStruct(storage.Ad{}, "json")

	propId, err := strconv.ParseUint(sPropId, 10, 64)
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

	prop, err := serviceProps.GetPropById(propId)
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
	prop.KindPropId = putRequest.KindPropId
	prop.Suffix = strings.TrimSpace(putRequest.Suffix)
	prop.Comment = strings.TrimSpace(putRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(putRequest.PrivateComment)

	if err = serviceProps.Update(prop, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	_, err = serviceProps.ReWriteValuesForProps(prop.PropId, tx, mId, mTitle, mPos)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	propFull, err := serviceProps.GetPropFullById(prop.PropId, serviceValuesProps)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = propFull
	return res
}
func deletePropsPropId(sPropId string) response.Result {
	serviceProps := service.NewPropService()
	res := response.Result{}

	propId, err := strconv.ParseUint(sPropId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()
	if err := serviceProps.Delete(propId, tx); err != nil {
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
