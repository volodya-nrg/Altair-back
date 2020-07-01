package controller

import (
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetMain - получение начальных данных, которые используются на фронте. Это начальные установки.
func GetMain(c *gin.Context) {
	roleIs := c.MustGet("roleIs").(string)
	serviceCats := service.NewCatService()
	serviceKindProps := service.NewKindPropService()
	serviceProps := service.NewPropService()
	isDisabledCats := 0

	if roleIs == manager.IsAdmin {
		isDisabledCats = -1
	}

	cats, err := serviceCats.GetCats(isDisabledCats)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	catsTree := serviceCats.GetCatsAsTree(cats)

	kindProps, err := serviceKindProps.GetKindProps("kind_prop_id asc")
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	props, err := serviceProps.GetProps("title asc")
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"catsTree":  catsTree,
		"kindProps": kindProps,
		"props":     props,
	})
}

// private -------------------------------------------------------------------------------------------------------------
