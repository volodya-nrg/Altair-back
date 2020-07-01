package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

// GetProps - получение всех свойств
func GetProps(c *gin.Context) {
	sCatID := c.DefaultQuery("catID", "0")
	serviceProps := service.NewPropService()

	catID, err := manager.SToUint64(sCatID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// вытаскивать propFull слишком накладно (1.4мб), поэтому возьмем только для определенной категории
	if catID > 0 {
		propsFull, err := serviceProps.GetPropsFullByCatID(catID, false)
		if err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, propsFull)
		return
	}

	propsWithKindName, err := serviceProps.GetPropsWithKindName()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, propsWithKindName)
}

// GetPropsPropID - получение конкретного свойства
func GetPropsPropID(c *gin.Context) {
	sPropID := c.Param("propID")
	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()

	propID, err := manager.SToUint64(sPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	propFull, err := serviceProps.GetPropFullByID(propID, serviceValuesProps)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, propFull)
}

// PostProps - добавление свойства
func PostProps(c *gin.Context) {
	postRequest := new(request.PostProp)

	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()
	serviceKindProp := service.NewKindPropService()
	prop := new(storage.Prop)

	kindPropID, err := manager.SToUint64(postRequest.KindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	kindPropData, err := serviceKindProp.GetKindPropByID(kindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	tx := server.Db.Begin()

	prop.Title = strings.TrimSpace(postRequest.Title)
	prop.KindPropID = kindPropID
	prop.Name = strings.TrimSpace(postRequest.Name)
	prop.Suffix = strings.TrimSpace(postRequest.Suffix)
	prop.Comment = strings.TrimSpace(postRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(postRequest.PrivateComment)

	if err := serviceProps.Create(prop, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// если это select, radio, (checkbox) - то принимаем данные от values, иначе только одно значение
	if ok, _ := manager.InArray(kindPropData.Name, manager.TagKindNumber); ok {
		values := make([]storage.ValueProp, 0)

		// т.к. пришел из вне 0, вставим необходимый id св-ва (к кому он принадлежит)
		for _, v := range postRequest.Values {
			values = append(values, storage.ValueProp{
				Title:  v.Title,
				Pos:    v.Pos,
				PropID: prop.PropID,
			})
		}

		// если данные приходят пачкой, через textarea, то перезапишем их в соот-ие map-ы
		selectAsTextarea := strings.TrimSpace(postRequest.SelectAsTextarea)
		if selectAsTextarea != "" {
			values = createMapsFromMultiText(prop.PropID, selectAsTextarea)
		}

		_, err = serviceProps.ReWriteValuesForProps(prop.PropID, tx, values)
		if err != nil {
			tx.Rollback()
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}
	}

	tx.Commit()

	propFull, err := serviceProps.GetPropFullByID(prop.PropID, serviceValuesProps)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(201, propFull)
}

// PutPropsPropID - изменение свойства
func PutPropsPropID(c *gin.Context) {
	sPropID := c.Param("propID")
	putRequest := new(request.PutProp)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	serviceProps := service.NewPropService()
	serviceValuesProps := service.NewValuesPropService()
	serviceKindProp := service.NewKindPropService()

	propID, err := manager.SToUint64(sPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	prop, err := serviceProps.GetPropByID(propID)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	kindPropID, err := manager.SToUint64(putRequest.KindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	prop.Title = strings.TrimSpace(putRequest.Title)
	prop.Name = strings.TrimSpace(putRequest.Name)
	prop.KindPropID = kindPropID
	prop.Suffix = strings.TrimSpace(putRequest.Suffix)
	prop.Comment = strings.TrimSpace(putRequest.Comment)
	prop.PrivateComment = strings.TrimSpace(putRequest.PrivateComment)

	kindPropData, err := serviceKindProp.GetKindPropByID(kindPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	tx := server.Db.Begin()

	if err = serviceProps.Update(prop, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// если это select, radio, (checkbox) - то принимаем данные от values, иначе только одно значение
	if ok, _ := manager.InArray(kindPropData.Name, manager.TagKindNumber); ok {
		if _, err = serviceProps.ReWriteValuesForProps(prop.PropID, tx, putRequest.Values); err != nil {
			tx.Rollback()
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}
	}

	tx.Commit()

	propFull, err := serviceProps.GetPropFullByID(prop.PropID, serviceValuesProps)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, propFull)
}

// DeletePropsPropID - удаление конкретного свойства
func DeletePropsPropID(c *gin.Context) {
	sPropID := c.Param("propID")
	serviceProps := service.NewPropService()

	propID, err := manager.SToUint64(sPropID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	tx := server.Db.Begin()

	if err := serviceProps.Delete(propID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	c.JSON(204, nil)
}

// private -------------------------------------------------------------------------------------------------------------
func createMapsFromMultiText(propID uint64, multiText string) []storage.ValueProp {
	aStrings := strings.Split(multiText, "\n")
	values := make([]storage.ValueProp, 0)

	for k, v := range aStrings {
		tmpValue := storage.ValueProp{
			Title:  strings.TrimSpace(v),
			Pos:    uint64(k + 1),
			PropID: propID,
		}

		values = append(values, tmpValue)
	}

	return values
}
