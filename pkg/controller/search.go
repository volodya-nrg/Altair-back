package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"strings"
	"unicode/utf8"
)

// GetSearchAds - получение объявлений через "поиск"
func GetSearchAds(c *gin.Context) {
	req := new(request.GetSearchAds)

	if err := c.ShouldBind(req); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceAds := service.NewAdService()
	query := strings.TrimSpace(req.Query)

	if utf8.RuneCountInString(query) < 2 {
		c.JSON(400, manager.ErrQueryStringIsShort.Error())
		return
	}

	adsFull, err := serviceAds.GetAdsFullBySearchTitle(query, req.CatID, req.Limit, req.Offset, c.Request.Form)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, adsFull)
}

// private -------------------------------------------------------------------------------------------------------------
