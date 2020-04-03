package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mime/multipart"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func GetAds(c *gin.Context) {
	res := getAds(c.DefaultQuery("catId", ""))
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func GetAdsAdId(c *gin.Context) {
	res := getAdsAdId(c.Param("adId"))
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PostAds(c *gin.Context) {
	pPostRequest := new(request.PostAd)
	if err := c.ShouldBind(pPostRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err)
		c.JSON(500, err.Error())
		return
	}

	res := postAds(pPostRequest, form, c.SaveUploadedFile, c.Request.PostForm)
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutAdsAdId(c *gin.Context) {
	sAdId := c.Param("adId")
	pPutRequest := new(request.PutAd)

	if err := c.ShouldBind(pPutRequest); err != nil {
		logger.Warning.Println(err)
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err)
		c.JSON(500, err.Error())
		return
	}

	res := putAds(sAdId, pPutRequest, form, c.SaveUploadedFile, c.Request.PostForm)
	if res.Err != nil {
		logger.Warning.Println(err)
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func DeleteAdsAdId(c *gin.Context) {
	res := deleteAdsAdId(c.Param("adId"))
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getAds(catIdSrc string) response.Result {
	serviceAds := service.NewAdService()
	serviceCats := service.NewCatService()
	serviceImages := service.NewImageService()
	serviceAdDetails := service.NewAdDetailService()
	res := response.Result{}

	if catIdSrc == "" {
		catIdSrc = "0"
	}

	catId, err := strconv.ParseUint(catIdSrc, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	pCats, err := serviceCats.GetCats()
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}
	pCatsTree := serviceCats.GetCatsAsTree(pCats)

	pCatsDescendants := serviceCats.GetDescendantsNastedLoop(pCatsTree, catId)
	sliceCatIds := make([]uint64, 0)
	sliceCatIds = append(sliceCatIds, serviceCats.GetIdsFromCatsTree(pCatsDescendants)...)
	sort.Slice(sliceCatIds, func(i, j int) bool { return sliceCatIds[i] < sliceCatIds[j] })

	pAdsFull, err := serviceAds.GetAdsFull(sliceCatIds, false, true, serviceImages, serviceAdDetails)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pAdsFull
	return res
}
func getAdsAdId(sAdId string) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceAdDetails := service.NewAdDetailService()
	res := response.Result{}

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 400
		res.Err = err
		return res
	}

	adFull, err := serviceAds.GetAdFullById(adId, serviceImages, serviceAdDetails)
	if gorm.IsRecordNotFoundError(err) {
		logger.Warning.Println(err)
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		logger.Warning.Println(err)
		res.Status = 400
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = adFull
	return res
}
func postAds(pPostRequest *request.PostAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error, postForm url.Values) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceProperties := service.NewPropertyService()
	serviceAdDetails := service.NewAdDetailService()
	serviceValuesProperty := service.NewValuesPropertyService()
	res := response.Result{}
	ad := new(storage.Ad)
	tx := server.Db.Debug().Begin()

	ad.Title = strings.TrimSpace(pPostRequest.Title)
	ad.CatId = pPostRequest.CatId
	ad.UserId = pPostRequest.UserId
	ad.Price = pPostRequest.Price
	ad.Text = strings.TrimSpace(pPostRequest.Text)

	if err := serviceAds.Create(ad, tx); err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	propsFull, err := serviceProperties.GetPropertiesFullByCatId(ad.CatId, false, serviceValuesProperty)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	adDetails, err := serviceAdDetails.BuildDataFromRequestFormAndCatProps(ad.AdId, postForm, propsFull)
	if err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceAdDetails.Create(adDetails, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	for _, file := range form.File["files"] {
		newFilePath, err := helpers.UploadImage(file, "./web/images", fnUpload)
		if err != nil {
			logger.Warning.Println(err)
			continue
		}

		image := new(storage.Image)
		image.Filepath = newFilePath
		image.ElId = ad.AdId
		image.Opt = "ad"

		if err := serviceImages.Create(image); err != nil {
			res.Status = 500
			res.Err = err
			return res
		}
	}

	adFull, err := serviceAds.GetAdFullById(ad.AdId, serviceImages, serviceAdDetails)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = adFull
	return res
}
func putAds(sAdId string, pPutRequest *request.PutAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error, postForm url.Values) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceProperties := service.NewPropertyService()
	serviceAdDetails := service.NewAdDetailService()
	serviceValuesProperty := service.NewValuesPropertyService()
	res := response.Result{}

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	pAd, err := serviceAds.GetAdById(adId)
	if gorm.IsRecordNotFoundError(err) {
		res.Status = 404
		res.Err = err
		return res

	} else if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()

	pAd.Title = strings.TrimSpace(pPutRequest.Title)
	pAd.CatId = pPutRequest.CatId
	pAd.UserId = pPutRequest.UserId
	pAd.Price = pPutRequest.Price
	pAd.Text = strings.TrimSpace(pPutRequest.Text)
	pAd.IsDisabled = pPutRequest.IsDisabled

	if err = serviceAds.Update(pAd, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	// достанем св-ва данной категории
	propsFull, err := serviceProperties.GetPropertiesFullByCatId(pAd.CatId, false, serviceValuesProperty)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	pAdDetailsNew, err := serviceAdDetails.BuildDataFromRequestFormAndCatProps(pAd.AdId, postForm, propsFull)
	if err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceAdDetails.Update(pAd.AdId, pAdDetailsNew, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	// обработаем текущие фото. Если что-то удалили (на фронте), то и удалим на беке.
	images, err := serviceImages.GetImagesByElIdsAndOpt([]uint64{pAd.AdId}, "ad")
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	if len(images) != len(pPutRequest.FilesAlreadyHas) {
		diffImages := make([]*storage.Image, 0)

		for _, image := range images {
			var has bool
			for _, requestImageFile := range pPutRequest.FilesAlreadyHas {
				if requestImageFile == image.Filepath {
					has = true
					break
				}
			}
			if !has {
				diffImages = append(diffImages, image)
			}
		}

		if len(diffImages) > 0 {
			if err := serviceImages.DeleteAll(diffImages); err != nil {
				res.Status = 500
				res.Err = err
				return res
			}
		}
	}
	//\

	for _, file := range form.File["files"] {
		newFilePath, err := helpers.UploadImage(file, "./web/images", fnUpload)
		if err != nil {
			logger.Warning.Println(err)
			continue
		}

		image := new(storage.Image)
		image.Filepath = newFilePath
		image.ElId = pAd.AdId
		image.Opt = "ad"

		if err := serviceImages.Create(image); err != nil {
			res.Status = 500
			res.Err = err
			return res
		}
	}

	pAdFull, err := serviceAds.GetAdFullById(pAd.AdId, serviceImages, serviceAdDetails)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = pAdFull
	return res
}
func deleteAdsAdId(sAdId string) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	res := response.Result{}

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	if err := serviceAds.Delete(adId, nil, serviceImages); err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
