package controller

import (
	"altair/api/request"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
)

func GetAds(c *gin.Context) {
	pResult := getAds(c.DefaultQuery("catId", ""))
	if pResult.Err != nil {
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func GetAdsAdId(c *gin.Context) {
	pResult := getAdsAdId(c.Param("adId"))
	if pResult.Err != nil {
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
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

	pResult := postAds(pPostRequest, form, c.SaveUploadedFile)
	if pResult.Err != nil {
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
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

	pResult := putAds(sAdId, pPutRequest, form, c.SaveUploadedFile)
	if pResult.Err != nil {
		logger.Warning.Println(err)
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}
func DeleteAdsAdId(c *gin.Context) {
	pResult := deleteAdsAdId(c.Param("adId"))
	if pResult.Err != nil {
		pResult.Data = pResult.Err.Error()
	}

	c.JSON(pResult.Status, pResult.Data)
}

// private -------------------------------------------------------------------------------------------------------------
func getAds(catIdSrc string) *result {
	serviceAds := service.NewAdService()
	serviceCats := service.NewCatService()
	res := new(result)

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

	catsTree, err := serviceCats.GetCatsAsTree()
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}
	catsDescendants := serviceCats.GetDescendantsNastedLoop(catsTree, catId)

	sliceCatIds := make([]uint64, 0)
	sliceCatIds = append(sliceCatIds, catsDescendants.CatId)
	sliceCatIds = append(sliceCatIds, serviceCats.GetIdsFromCatsTree(catsDescendants)...)
	sort.Slice(sliceCatIds, func(i, j int) bool { return sliceCatIds[i] < sliceCatIds[j] })

	ads, err := serviceAds.GetAdsFull(sliceCatIds)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = ads
	return res
}
func getAdsAdId(sAdId string) *result {
	serviceAds := service.NewAdService()
	res := new(result)

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 400
		res.Err = err
		return res
	}

	adFull, err := serviceAds.GetAdFullById(adId)
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
func postAds(pPostRequest *request.PostAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) *result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	res := new(result)
	ad := new(storage.Ad)

	ad.Title = strings.TrimSpace(pPostRequest.Title)
	ad.CatId = pPostRequest.CatId
	ad.UserId = pPostRequest.UserId
	ad.Price = pPostRequest.Price
	ad.Text = strings.TrimSpace(pPostRequest.Text)

	if err := serviceAds.Create(ad); err != nil {
		res.Status = 400
		res.Err = err
		return res
	}

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

	adFull, err := serviceAds.GetAdFullById(ad.AdId)
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
func putAds(sAdId string, pPutRequest *request.PutAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) *result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	res := new(result)

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

	pAd.Title = strings.TrimSpace(pPutRequest.Title)
	pAd.CatId = pPutRequest.CatId
	pAd.UserId = pPutRequest.UserId
	pAd.Price = pPutRequest.Price
	pAd.Text = strings.TrimSpace(pPutRequest.Text)
	pAd.IsDisabled = pPutRequest.IsDisabled

	if err = serviceAds.Update(pAd); err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	// обработаем текущие фото. Если что-то удалили (на фронте), то и удалим на беке.
	images, err := serviceImages.GetImages(pAd.AdId, "ad")
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

	pAdFull, err := serviceAds.GetAdFullById(pAd.AdId)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 201
	res.Err = nil
	res.Data = pAdFull
	return res
}
func deleteAdsAdId(sAdId string) *result {
	serviceAds := service.NewAdService()
	res := new(result)

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	if err := serviceAds.Delete(adId); err != nil {
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

type result struct {
	Status int
	Err    error
	Data   interface{}
}
