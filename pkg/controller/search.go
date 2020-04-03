package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/logger"
	"altair/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
	"unicode/utf8"
)

func GetSearchAds(c *gin.Context) {
	getRequest := request.GetSearchAds{}

	if err := c.ShouldBind(&getRequest); err != nil {
		c.JSON(400, err.Error())
		return
	}

	res := getSearchAds(getRequest.Query, getRequest.CatId, c.Request.Form)
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getSearchAds(querySrc string, catId uint64, mGetParams url.Values) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceAdDetails := service.NewAdDetailService()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}
	query := strings.TrimSpace(querySrc)

	if utf8.RuneCountInString(query) < 2 {
		logger.Warning.Println("query string is short")
		res.Status = 400
		res.Err = errors.New("query string is short")
		return res
	}

	pAdsFull, err := serviceAds.GetAdsFullBySearchTitle(query, catId, mGetParams, serviceImages, serviceAdDetails, serviceProperties, serviceValuesProperties)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pAdsFull
	return res
}
