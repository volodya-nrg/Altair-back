package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// GetCats - получение всех категорий
func GetCats(c *gin.Context) {
	asTree := c.DefaultQuery("asTree", "false")
	serviceCats := service.NewCatService()
	isAsTree := asTree == "true"

	cats, err := serviceCats.GetCats(0)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	if isAsTree {
		c.JSON(200, serviceCats.GetCatsAsTree(cats))
		return
	}

	c.JSON(200, cats)
}

// GetCatsCatID - получение одной категории
func GetCatsCatID(c *gin.Context) {
	roleIs := c.MustGet("roleIs").(string)
	sCatID := c.Param("catID")
	sWithPropsOnlyFiltered := c.DefaultQuery("withPropsOnlyFiltered", "false")
	serviceCats := service.NewCatService()
	isDisabledCat := 0

	if roleIs == manager.IsAdmin {
		isDisabledCat = -1
	}

	catID, err := manager.SToUint64(sCatID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	catFull, err := serviceCats.GetCatFullByID(catID, sWithPropsOnlyFiltered == "true", isDisabledCat)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, catFull)
}

// PostCats - создание одной категории
func PostCats(c *gin.Context) {
	postRequest := new(request.PostCat)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceCats := service.NewCatService()

	parentID, err := manager.SToUint64(postRequest.ParentID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	cat := new(storage.Cat)
	cat.Name = strings.TrimSpace(postRequest.Name)
	cat.ParentID = parentID
	cat.Pos = postRequest.Pos
	cat.PriceAlias = strings.TrimSpace(postRequest.PriceAlias)
	cat.PriceSuffix = strings.TrimSpace(postRequest.PriceSuffix)
	cat.TitleHelp = strings.TrimSpace(postRequest.TitleHelp)
	cat.TitleComment = strings.TrimSpace(postRequest.TitleComment)
	cat.IsAutogenerateTitle = postRequest.IsAutogenerateTitle

	if cat.Pos < 1 {
		cat.Pos = 1
	}

	tx := server.Db.Begin()

	if err := serviceCats.Create(cat, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// обработаем св-ва для категории
	if _, err := serviceCats.ReWriteCatsProps(cat.CatID, tx, postRequest.PropsAssignedForCat); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	catFull, err := serviceCats.GetCatFullByID(cat.CatID, false, -1)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(201, catFull)
}

// PutCatsCatID - редактирование одной категории
func PutCatsCatID(c *gin.Context) {
	sCatID := c.Param("catID")
	putRequest := new(request.PutCat)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceCats := service.NewCatService()

	catID, err := manager.SToUint64(sCatID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	cat, err := serviceCats.GetCatByID(catID, -1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	tx := server.Db.Begin()

	cat.Name = strings.TrimSpace(putRequest.Name)
	cat.ParentID = putRequest.ParentID
	cat.Pos = putRequest.Pos
	cat.IsDisabled = putRequest.IsDisabled
	cat.PriceAlias = strings.TrimSpace(putRequest.PriceAlias)
	cat.PriceSuffix = strings.TrimSpace(putRequest.PriceSuffix)
	cat.TitleHelp = strings.TrimSpace(putRequest.TitleHelp)
	cat.TitleComment = strings.TrimSpace(putRequest.TitleComment)
	cat.IsAutogenerateTitle = putRequest.IsAutogenerateTitle

	if cat.Pos < 1 {
		cat.Pos = 1
	}

	if err = serviceCats.Update(cat, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	if _, err := serviceCats.ReWriteCatsProps(cat.CatID, tx, putRequest.PropsAssignedForCat); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	catFull, err := serviceCats.GetCatFullByID(cat.CatID, false, -1)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, catFull)
}

// DeleteCatsCatID - удаление одной категории
func DeleteCatsCatID(c *gin.Context) {
	sCatID := c.Param("catID")
	serviceCats := service.NewCatService()

	catID, err := manager.SToUint64(sCatID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	tx := server.Db.Begin()

	if err := serviceCats.Delete(catID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	c.JSON(204, nil)
}

// private -------------------------------------------------------------------------------------------------------------
