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

func GetKindProps(c *gin.Context) {
	res := getKindProps()
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetKindPropsKindPropId(c *gin.Context) {
	res := getKindPropsKindPropId(c.Param("kindPropId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostKindProps(c *gin.Context) {
	pPostRequest := new(request.PostKindProp)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := postKindProps(pPostRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutKindPropsKindPropId(c *gin.Context) {
	pPutRequest := new(request.PutKindProp)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	res := putKindPropsKindPropId(c.Param("kindPropId"), pPutRequest)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeleteKindPropsKindPropId(c *gin.Context) {
	res := deleteKindPropsKindPropId(c.Param("kindPropId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getKindProps() response.Result {
	serviceKindProps := service.NewKindPropService()
	res := response.Result{}

	pKindProps, err := serviceKindProps.GetKindProps("kind_prop_id desc")
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pKindProps
	return res
}
func getKindPropsKindPropId(sKindPropId string) response.Result {
	serviceKindProps := service.NewKindPropService()
	res := response.Result{}

	kindPropId, err := strconv.ParseUint(sKindPropId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	pKindProp, err := serviceKindProps.GetKindPropById(kindPropId)
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
	res.Data = pKindProp
	return res
}
func postKindProps(pPostRequest *request.PostKindProp) response.Result {
	serviceKindProps := service.NewKindPropService()
	res := response.Result{}
	pKindProp := new(storage.KindProp)

	pKindProp.Name = pPostRequest.Name

	if err := serviceKindProps.Create(pKindProp, nil); err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = pKindProp
	return res
}
func putKindPropsKindPropId(sKindPropId string, putRequest *request.PutKindProp) response.Result {
	serviceKindProps := service.NewKindPropService()
	res := response.Result{}

	kindPropId, err := strconv.ParseUint(sKindPropId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pKindProp, err := serviceKindProps.GetKindPropById(kindPropId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	pKindProp.Name = strings.TrimSpace(putRequest.Name)

	if err = serviceKindProps.Update(pKindProp, nil); err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pKindProp
	return res
}
func deleteKindPropsKindPropId(sKindPropId string) response.Result {
	serviceKindProps := service.NewKindPropService()
	res := response.Result{}

	kindPropId, err := strconv.ParseUint(sKindPropId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceKindProps.Delete(kindPropId, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
