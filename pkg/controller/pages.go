package controller

import (
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPagesAdAdID - получение данных для страницы "продукта"
func GetPagesAdAdID(c *gin.Context) {
	sAdID := c.Param("adID")
	serviceAds := service.NewAdService()
	serviceCats := service.NewCatService()

	adID, err := manager.SToUint64(sAdID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	adFull, err := serviceAds.GetAdFullByID(adID, 0, 1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warning.Println(err.Error())
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		c.JSON(500, err.Error())
		return
	}

	catFull, err := serviceCats.GetCatFullByID(adFull.CatID, false, 0)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"adFull":  adFull,
		"catFull": catFull,
	})
}

// GetPagesMain - получение данных для гл. страницы
func GetPagesMain(c *gin.Context) {
	sLimit := c.DefaultQuery("limit", "0")
	serviceAds := service.NewAdService()
	limit := manager.LimitDefault

	if limitTmp, err := manager.SToUint64(sLimit); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return

	} else if limitTmp > 0 && int(limitTmp) < limit {
		limit = int(limitTmp)
	}

	ads, err := serviceAds.GetLastAdsByOneCat(limit)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"lastAdsFull": ads,
	})
}

// private -------------------------------------------------------------------------------------------------------------
