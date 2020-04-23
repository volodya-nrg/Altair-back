package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetCats(c *gin.Context) {
	asTree := c.DefaultQuery("asTree", "false")
	res := response.Result{}

	res = getCats(asTree == "true")
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetCatsCatId(c *gin.Context) {
	withPropsOnlyFiltered := c.DefaultQuery("withPropsOnlyFiltered", "false")

	res := getCatsCatId(c.Param("catId"), withPropsOnlyFiltered == "true")
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostCats(c *gin.Context) {
	pPostRequest := new(request.PostCat)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	mPropertyId := c.PostFormMap("propertyId")
	mPos := c.PostFormMap("pos")
	mIsRequire := c.PostFormMap("isRequire")
	mIsCanAsFilter := c.PostFormMap("isCanAsFilter")
	mComment := c.PostFormMap("comment")

	res := postCats(pPostRequest, mPropertyId, mPos, mIsRequire, mIsCanAsFilter, mComment)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutCatsCatId(c *gin.Context) {
	pPutRequest := new(request.PutCat)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	mPropertyId := c.PostFormMap("propertyId")
	mPos := c.PostFormMap("pos")
	mIsRequire := c.PostFormMap("isRequire")
	mIsCanAsFilter := c.PostFormMap("isCanAsFilter")
	mComment := c.PostFormMap("comment")

	res := putCatsCatId(c.Param("catId"), pPutRequest, mPropertyId, mPos, mIsRequire, mIsCanAsFilter, mComment)
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeleteCatsCatId(c *gin.Context) {
	res := deleteCatsCatId(c.Param("catId"))
	if res.Err != nil {
		logger.Warning.Println(res.Err.Error())
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getCats(isAsTree bool) response.Result {
	serviceCats := service.NewCatService()
	res := response.Result{}

	cats, err := serviceCats.GetCats()
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = cats

	if isAsTree {
		res.Data = serviceCats.GetCatsAsTree(cats)
	}

	return res
}
func getCatsCatId(sCatId string, withPropsOnlyFiltered bool) response.Result {
	serviceCats := service.NewCatService()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	catFull, err := serviceCats.GetCatFullByID(catId, withPropsOnlyFiltered, serviceProperties, serviceValuesProperties)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = catFull
	return res
}
func postCats(postRequest *request.PostCat, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string, mIsCanAsFilter map[string]string, mComment map[string]string) response.Result {
	serviceCats := service.NewCatService()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}
	pCat := new(storage.Cat)
	tx := server.Db.Debug().Begin()

	pCat.Name = strings.TrimSpace(postRequest.Name)
	pCat.ParentId = postRequest.ParentId
	pCat.Pos = postRequest.Pos
	pCat.PriceAlias = strings.TrimSpace(postRequest.PriceAlias)
	pCat.PriceSuffix = strings.TrimSpace(postRequest.PriceSuffix)
	pCat.TitleHelp = strings.TrimSpace(postRequest.TitleHelp)
	pCat.TitleComment = strings.TrimSpace(postRequest.TitleComment)
	pCat.IsAutogenerateTitle = postRequest.IsAutogenerateTitle

	if pCat.Pos < 1 {
		pCat.Pos = 1
	}

	if err := serviceCats.Create(pCat, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	// обработаем св-ва для категории
	if _, err := serviceCats.ReWriteCatsProperties(pCat.CatId, tx, mPropertyId, mPos, mIsRequire, mIsCanAsFilter, mComment); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	catFull, err := serviceCats.GetCatFullByID(pCat.CatId, false, serviceProperties, serviceValuesProperties)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = catFull
	return res
}
func putCatsCatId(sCatId string, putRequest *request.PutCat, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string, mIsCanAsFilter map[string]string, mComment map[string]string) response.Result {
	serviceCats := service.NewCatService()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	res := response.Result{}

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pCat, err := serviceCats.GetCatByID(catId)
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

	pCat.Name = strings.TrimSpace(putRequest.Name)
	pCat.ParentId = putRequest.ParentId
	pCat.Pos = putRequest.Pos
	pCat.IsDisabled = putRequest.IsDisabled
	pCat.PriceAlias = strings.TrimSpace(putRequest.PriceAlias)
	pCat.PriceSuffix = strings.TrimSpace(putRequest.PriceSuffix)
	pCat.TitleHelp = strings.TrimSpace(putRequest.TitleHelp)
	pCat.TitleComment = strings.TrimSpace(putRequest.TitleComment)
	pCat.IsAutogenerateTitle = putRequest.IsAutogenerateTitle

	if pCat.Pos < 1 {
		pCat.Pos = 1
	}

	if err = serviceCats.Update(pCat, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	if _, err := serviceCats.ReWriteCatsProperties(pCat.CatId, tx, mPropertyId, mPos, mIsRequire, mIsCanAsFilter, mComment); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	catFull, err := serviceCats.GetCatFullByID(pCat.CatId, false, serviceProperties, serviceValuesProperties)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = catFull
	return res
}
func deleteCatsCatId(sCatId string) response.Result {
	serviceCats := service.NewCatService()
	res := response.Result{}

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceCats.Delete(catId, nil); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
