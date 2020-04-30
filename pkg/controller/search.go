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
	req := request.GetSearchAds{}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, err.Error())
		return
	}

	res := getSearchAds(req.Query, req.CatId, req.Limit, req.Offset, c.Request.Form)
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getSearchAds(querySrc string, catId uint64, limit uint64, offset uint64, mGetParams url.Values) response.Result {
	serviceAds := service.NewAdService()
	res := response.Result{}
	query := strings.TrimSpace(querySrc)

	if utf8.RuneCountInString(query) < 2 {
		logger.Warning.Println("query string is short")
		res.Status = 400
		res.Err = errors.New("query string is short")
		return res
	}

	pAdsFull, err := serviceAds.GetAdsFullBySearchTitle(query, catId, limit, offset, mGetParams)
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
