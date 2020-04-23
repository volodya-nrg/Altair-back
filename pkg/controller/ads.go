package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"fmt"
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
	postRequest := new(request.PostAd)
	if err := c.ShouldBind(postRequest); err != nil {
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

	res := postAds(postRequest, form, c.SaveUploadedFile, &c.Request.PostForm)
	if res.Err != nil {
		res.Data = res.Err.Error()
	}

	c.JSON(res.Status, res.Data)
}
func PutAdsAdId(c *gin.Context) {
	sAdId := c.Param("adId")
	putRequest := new(request.PutAd)

	if err := c.ShouldBind(putRequest); err != nil {
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

	res := putAds(sAdId, putRequest, form, c.SaveUploadedFile, &c.Request.PostForm)
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
	sliceCatIds := make([]uint64, 0)

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

	if catId > 0 {
		pCats, err := serviceCats.GetCats()
		if err != nil {
			logger.Warning.Println(err)
			res.Status = 500
			res.Err = err
			return res
		}

		pCatsTree := serviceCats.GetCatsAsTree(pCats)
		pCatsDescendants := serviceCats.GetDescendantsNastedLoop(pCatsTree, catId)
		sliceCatIds = append(sliceCatIds, serviceCats.GetIdsFromCatsTree(pCatsDescendants)...)
		sort.Slice(sliceCatIds, func(i, j int) bool { return sliceCatIds[i] < sliceCatIds[j] })
	}

	adsFull, err := serviceAds.GetAdsFull(sliceCatIds, false, true, serviceImages, serviceAdDetails)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = adsFull
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
func postAds(postRequest *request.PostAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error, postForm *url.Values) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceProperties := service.NewPropertyService()
	serviceAdDetails := service.NewAdDetailService()
	serviceValuesProperty := service.NewValuesPropertyService()
	res := response.Result{}
	ad := new(storage.Ad)
	tx := server.Db.Debug().Begin()

	ad.Title = strings.TrimSpace(postRequest.Title)
	ad.CatId = postRequest.CatId
	ad.UserId = postRequest.UserId
	ad.Price = postRequest.Price
	ad.Description = strings.TrimSpace(postRequest.Description)
	ad.Youtube = strings.TrimSpace(postRequest.Youtube)

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

	images := make([]*storage.Image, 0)
	if err := workWithPhoto(ad, images, propsFull, tx, form, serviceImages, serviceAds, postForm, fnUpload); err != nil {
		tx.Rollback()
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	// тут происходит просто валидная сборка приходящих данных с тем что должно быть
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
func putAds(sAdId string, putRequest *request.PutAd, form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error, postForm *url.Values) response.Result {
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

	ad, err := serviceAds.GetAdById(adId)
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

	ad.Title = strings.TrimSpace(putRequest.Title)
	ad.CatId = putRequest.CatId
	ad.UserId = putRequest.UserId
	ad.Price = putRequest.Price
	ad.Description = strings.TrimSpace(putRequest.Description)
	ad.IsDisabled = putRequest.IsDisabled
	ad.Youtube = strings.TrimSpace(putRequest.Youtube)

	if err = serviceAds.Update(ad, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	// достанем св-ва данной категории
	propsFull, err := serviceProperties.GetPropertiesFullByCatId(ad.CatId, false, serviceValuesProperty)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	// возьмем текущие фото
	images, err := serviceImages.GetImagesByElIdsAndOpt([]uint64{ad.AdId}, "ad")
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	// если есть разница, то удалим не нужное
	if len(images) != len(putRequest.FilesAlreadyHas) {
		restOfImages := make([]*storage.Image, 0) // срез остатков фото

		for _, image := range images {
			var has bool
			for _, requestImageFile := range putRequest.FilesAlreadyHas {
				if requestImageFile == image.Filepath {
					restOfImages = append(restOfImages, image)
					has = true
					break
				}
			}
			if !has {
				if err := serviceImages.Delete(image, tx); err != nil {
					tx.Rollback()
					res.Status = 500
					res.Err = err
					return res
				}
			}
		}

		images = restOfImages
	}

	if err := workWithPhoto(ad, images, propsFull, tx, form, serviceImages, serviceAds, postForm, fnUpload); err != nil {
		tx.Rollback()
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	adDetailsNew, err := serviceAdDetails.BuildDataFromRequestFormAndCatProps(ad.AdId, postForm, propsFull)
	if err != nil {
		tx.Rollback()
		res.Status = 400
		res.Err = err
		return res
	}

	if err := serviceAdDetails.Update(ad.AdId, adDetailsNew, tx); err != nil {
		tx.Rollback()
		res.Status = 500
		res.Err = err
		return res
	}

	tx.Commit()

	adFull, err := serviceAds.GetAdFullById(ad.AdId, serviceImages, serviceAdDetails)
	if err != nil {
		res.Status = 500
		res.Err = err
		return res
	}

	res.Status = 200
	res.Err = nil
	res.Data = adFull
	return res
}
func deleteAdsAdId(sAdId string) response.Result {
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceAdDetail := service.NewAdDetailService()
	res := response.Result{}

	adId, err := strconv.ParseUint(sAdId, 10, 64)
	if err != nil {
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}

	tx := server.Db.Debug().Begin()
	if err := serviceAds.Delete(adId, tx, serviceImages, serviceAdDetail); err != nil {
		tx.Rollback()
		logger.Warning.Println(err)
		res.Status = 500
		res.Err = err
		return res
	}
	tx.Commit()

	res.Status = 204
	res.Err = nil
	res.Data = nil
	return res
}
func workWithPhoto(ad *storage.Ad, curImages []*storage.Image, propsFull []*response.PropertyFull, tx *gorm.DB, form *multipart.Form, serviceImages *service.ImageService, serviceAds *service.AdService, postForm *url.Values, fnUpload func(file *multipart.FileHeader, filePath string) error) error {
	sliceImageSIds := make([]string, 0)

	// получить текущие фото
	for _, v := range curImages {
		sliceImageSIds = append(sliceImageSIds, fmt.Sprint(v.ImgId))
	}

	for _, propFull := range propsFull {
		if propFull.KindPropertyName != "photo" { // вид св-ва - photo
			continue
		}

		max, err := strconv.Atoi(propFull.Comment)
		if err != nil {
			return err
		}

		max = max - len(curImages)

		for key, file := range form.File["files"] {
			// обойдем только определенное число файлов
			if key >= max {
				break
			}

			newFilePath, err := helpers.UploadImage(file, "./web/images", fnUpload)
			if err != nil {
				logger.Warning.Println(err)
				continue
			}

			image := new(storage.Image)
			image.Filepath = newFilePath
			image.ElId = ad.AdId
			image.Opt = "ad"

			if err := serviceImages.Create(image, tx); err != nil {
				return err
			}

			sliceImageSIds = append(sliceImageSIds, fmt.Sprint(image.ImgId))
		}

		// если есть что добавлять, добавим
		if len(sliceImageSIds) > 0 {
			postForm.Del("files")
			postForm.Add("files", strings.Join(sliceImageSIds, ",")) // добавляем POST переменную для фото

			// обновим данные у объявления о наличии фото, чтоб удобно потом считать
			ad.HasPhoto = true
			if err := serviceAds.Update(ad, tx); err != nil {
				return err
			}
		}

		break // только один раз
	}

	return nil
}
