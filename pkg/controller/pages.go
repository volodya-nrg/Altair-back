package controller

import (
	"altair/api/response"
	"altair/pkg/logger"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetPagesAdAdId(c *gin.Context) {
	res, status, err := getPagesAdAdId(c.Param("adId"))
	if err != nil {
		c.JSON(status, err.Error())
		return
	}

	c.JSON(status, res)
}
func GetPagesMain(c *gin.Context) {
	limitSrc := c.DefaultQuery("limit", "0")

	res, status, err := getPagesMain(limitSrc)
	if err != nil {
		c.JSON(status, err.Error())
		return
	}

	c.JSON(status, res)
}

// private -------------------------------------------------------------------------------------------------------------
func getPagesAdAdId(sAdId string) (*response.PageAd, int, error) {
	serviceAds := service.NewAdService()
	serviceCats := service.NewCatService()
	pageAd := new(response.PageAd)

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		return pageAd, 400, err
	}

	adFull, err := serviceAds.GetAdFullById(adId)
	if gorm.IsRecordNotFoundError(err) {
		logger.Warning.Println(err)
		return pageAd, 404, nil

	} else if err != nil {
		logger.Warning.Println(err)
		return pageAd, 500, err
	}

	pageAd.AdFull = adFull

	catFull, err := serviceCats.GetCatFullByID(adFull.CatId, false)
	if gorm.IsRecordNotFoundError(err) {
		logger.Warning.Println(err)
		return pageAd, 404, nil

	} else if err != nil {
		logger.Warning.Println(err)
		return pageAd, 500, err
	}

	pageAd.CatFull = catFull

	return pageAd, 200, nil
}
func getPagesMain(sLimit string) (*response.PageMain, int, error) {
	serviceAds := service.NewAdService()
	pageMain := new(response.PageMain)
	var limit uint64

	if tmpLimit, err := strconv.ParseUint(sLimit, 10, 64); err == nil {
		limit = tmpLimit
	}

	if limit < 1 || limit > 10 {
		limit = 4
	}

	ads, err := serviceAds.GetLastAdsByOneCat(limit)
	if err != nil {
		return pageMain, 500, err
	}

	pageMain.Last.AdsFull = ads

	return pageMain, 200, nil
}
