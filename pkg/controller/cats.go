package controller

import (
	"altair/api/request"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetCats(c *gin.Context) {
	asTree := c.DefaultQuery("asTree", "false")
	isFillPropertiesFull := c.DefaultQuery("isFillPropertiesFull", "false")
	pResult := new(result)

	if isFillPropertiesFull == "true" {
		pResult = getCatsFull(asTree == "true")
		if pResult.Err != nil {
			logger.Warning.Println(pResult.Err.Error())
			pResult.Data = pResult.Err.Error()
		}

	} else {
		pResult = getCats(asTree == "true")
		if pResult.Err != nil {
			logger.Warning.Println(pResult.Err.Error())
			pResult.Data = pResult.Err.Error()
		}
	}

	c.JSON(pResult.Status, pResult.Data)
}
func GetCatsCatId(c *gin.Context) {
	pResult := getCatsCatId(c.Param("catId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PostCats(c *gin.Context) {
	pPostRequest := new(request.PostCat)

	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := postCats(pPostRequest, c.PostFormMap("propertyId"), c.PostFormMap("pos"), c.PostFormMap("isRequire"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func PutCatsCatId(c *gin.Context) {
	pPutRequest := new(request.PutCat)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	pResult := putCatsCatId(c.Param("catId"), pPutRequest, c.PostFormMap("propertyId"), c.PostFormMap("pos"), c.PostFormMap("isRequire"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func DeleteCatsCatId(c *gin.Context) {
	pResult := deleteCatsCatId(c.Param("catId"))
	if pResult.Err != nil {
		logger.Warning.Println(pResult.Err.Error())
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getCats(isAsTree bool) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)

	cats, err := serviceCats.GetCats()
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = cats

	if isAsTree {
		pResult.Data = serviceCats.GetCatsAsTree(cats)
	}

	return pResult
}
func getCatsFull(isAsTree bool) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)

	cats, err := serviceCats.GetCatsFull()
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = cats

	if isAsTree {
		pResult.Data = serviceCats.GetCatsFullAsTree(cats)
	}

	return pResult
}
func getCatsCatId(sCatId string) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	catFull, err := serviceCats.GetCatFullByID(catId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = catFull
	return pResult
}
func postCats(postRequest *request.PostCat, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)
	pCat := new(storage.Cat)

	pCat.Name = postRequest.Name
	pCat.ParentId = postRequest.ParentId
	pCat.Pos = postRequest.Pos

	if err := serviceCats.Create(pCat); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	// обработаем св-ва для категории
	if _, err := serviceCats.ReWriteCatsProperties(pCat.CatId, mPropertyId, mPos, mIsRequire); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	catFull, err := serviceCats.GetCatFullByID(pCat.CatId)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 201
	pResult.Err = nil
	pResult.Data = catFull
	return pResult
}
func putCatsCatId(sCatId string, putRequest *request.PutCat, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pCat, err := serviceCats.GetCatByID(catId)
	if gorm.IsRecordNotFoundError(err) {
		pResult.Status = 404
		pResult.Err = err
		return pResult

	} else if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	pCat.Name = strings.TrimSpace(putRequest.Name)
	pCat.ParentId = putRequest.ParentId
	pCat.Pos = putRequest.Pos
	pCat.IsDisabled = putRequest.IsDisabled

	if err = serviceCats.Update(pCat); err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	if _, err := serviceCats.ReWriteCatsProperties(pCat.CatId, mPropertyId, mPos, mIsRequire); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	catFull, err := serviceCats.GetCatFullByID(pCat.CatId)
	if err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 200
	pResult.Err = nil
	pResult.Data = catFull
	return pResult
}
func deleteCatsCatId(sCatId string) *result {
	serviceCats := service.NewCatService()
	pResult := new(result)

	catId, err := strconv.ParseUint(sCatId, 10, 64)
	if err != nil {
		pResult.Status = 400
		pResult.Err = err
		return pResult
	}

	if err := serviceCats.Delete(catId); err != nil {
		pResult.Status = 500
		pResult.Err = err
		return pResult
	}

	pResult.Status = 204
	pResult.Err = nil
	pResult.Data = nil
	return pResult
}
