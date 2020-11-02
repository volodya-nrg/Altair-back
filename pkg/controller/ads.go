package controller

import (
	"altair/api/request"
	"altair/api/response"
	"altair/configs"
	"altair/pkg/emailer"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/mediafire"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

// GetAds - получить все объявления
func GetAds(c *gin.Context) {
	roleIs := c.MustGet("roleIs").(string)
	sCatIDSrc := c.DefaultQuery("catID", "0")
	sLimitSrc := c.DefaultQuery("limit", "0")
	sOffsetSrc := c.DefaultQuery("offset", "0")
	serviceAds := service.NewAdService()
	serviceCats := service.NewCatService()
	sliceCatIDs := make([]uint64, 0)
	limit := manager.LimitDefault
	isDisabledCats := 0
	isDisabledAds := 0
	isApprovedAds := 1
	var catID uint64
	var offset uint64

	if roleIs == manager.IsAdmin {
		isDisabledCats = -1
		isDisabledAds = -1
		isApprovedAds = -1
	}

	if catIDTmp, err := manager.SToUint64(sCatIDSrc); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	} else if catIDTmp > 0 {
		catID = catIDTmp
	}

	if limitTmp, err := manager.SToUint64(sLimitSrc); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	} else if limitTmp > 0 && int(limitTmp) < manager.LimitDefault {
		limit = int(limitTmp)
	}

	if offsetTmp, err := manager.SToUint64(sOffsetSrc); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	} else if offsetTmp > 0 {
		offset = offsetTmp
	}

	if catID > 0 {
		cats, err := serviceCats.GetCats(isDisabledCats)
		if err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		}

		catsTree := serviceCats.GetCatsAsTree(cats)
		catsDescendants := serviceCats.GetDescendants(catsTree, catID) // потомки
		sliceCatIDs = append(sliceCatIDs, serviceCats.GetIDsFromCatsTree(catsDescendants)...)
		sort.Slice(sliceCatIDs, func(i, j int) bool { return sliceCatIDs[i] < sliceCatIDs[j] })
	}

	adsFull, err := serviceAds.GetAdsFull(
		sliceCatIDs,
		false,
		"created_at desc",
		limit,
		offset,
		isDisabledAds,
		isApprovedAds)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, adsFull)
}

// GetAdsAdID - получить объявление по id
func GetAdsAdID(c *gin.Context) {
	roleIs := c.MustGet("roleIs").(string)
	sAdID := c.Param("adID")
	serviceAds := service.NewAdService()
	isDisabledAd := 0
	isApprovedAd := 1

	if roleIs == manager.IsAdmin {
		isDisabledAd = -1
		isApprovedAd = -1
	}

	adID, err := manager.SToUint64(sAdID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	adFull, err := serviceAds.GetAdFullByID(adID, isDisabledAd, isApprovedAd)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, adFull)
}

// PostAds - создание объявления
func PostAds(c *gin.Context) {
	postRequest := new(request.PostAd)
	if err := c.ShouldBind(postRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserID.Error())
		return
	}

	userRole, ok := c.MustGet("userRole").(string)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserRole.Error())
		return
	}
	//----------------------------------------------------------

	serviceAds := service.NewAdService()
	serviceAdDetails := service.NewAdDetailService()
	serviceCat := service.NewCatService()
	isDisabledCat := 0
	isDisabledAd := 0
	isApprovedAd := 0
	ad := new(storage.Ad)

	ad.Title = strings.TrimSpace(postRequest.Title)
	ad.CatID = postRequest.CatID
	ad.Price = postRequest.Price
	ad.Description = strings.TrimRight(postRequest.Description, " ")
	ad.Latitude = postRequest.Latitude
	ad.Longitude = postRequest.Longitude
	ad.CityName = strings.TrimSpace(postRequest.CityName)
	ad.CountryName = strings.TrimSpace(postRequest.CountryName)
	ad.IP = c.ClientIP()
	ad.Youtube = manager.GetYoutubeHash(postRequest.Youtube)
	ad.PhoneID = postRequest.PhoneID

	if userRole == manager.IsAdmin {
		isDisabledCat = -1
		isDisabledAd = -1
		isApprovedAd = -1

		ad.UserID = postRequest.UserID // вот тут может быть неск-ко вариантов
		ad.IsDisabled = postRequest.IsDisabled
		ad.IsApproved = postRequest.IsApproved

		if postRequest.UserID == 0 {
			ad.UserID = userID
		}

	} else {
		ad.UserID = userID
	}

	if statusCode, err := checkPhoneNumber(ad); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(statusCode, err.Error())
		return
	}

	catFull, err := serviceCat.GetCatFullByID(ad.CatID, false, isDisabledCat)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// первая проверка объявления (на корректность)
	if err := firstCheck(ad, catFull.Cat); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// если автогинерация, то подменим. Заголовок обязательно должен быть.
	if catFull.IsAutogenerateTitle {
		updatedTitle, err := updateTitle(catFull, &c.Request.PostForm)
		if err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(400, err.Error())
			return

		} else if updatedTitle == "" {
			c.JSON(400, manager.ErrTitleIsEmpty.Error())
			return
		}

		ad.Title = updatedTitle
	}

	tx := server.Db.Begin()

	if err := serviceAds.Create(ad, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	images := make([]*storage.Image, 0)
	if err := workWithPhoto(ad, images, catFull.PropsFull, tx, form, c.SaveUploadedFile); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// тут происходит просто валидная сборка приходящих данных с тем что должно быть
	adDetails, err := serviceAdDetails.BuildDataFromRequestFormAndCatProps(ad.AdID, &c.Request.PostForm, catFull.PropsFull, images)
	if err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	if err := serviceAdDetails.Create(adDetails, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	adFull, err := serviceAds.GetAdFullByID(ad.AdID, isDisabledAd, isApprovedAd)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	var prefixTest string
	if configs.Cfg.Mode == "debug" {
		prefixTest = ", test"
	}

	emailRequest := emailer.NewEmailRequest(
		configs.Cfg.AdminEmail,
		"Altair (add"+prefixTest+"): "+adFull.Title,
		adFull.Description)

	if ok, err := emailRequest.SendMail(); err != nil {
		logger.Warning.Println(err.Error())

	} else if !ok {
		logger.Warning.Println(manager.ErrEmailNotSend.Error())
	}

	c.JSON(201, adFull)
}

// PutAdsAdID - редактирование объявления
func PutAdsAdID(c *gin.Context) {
	sAdID := c.Param("adID")
	putRequest := new(request.PutAd)

	if err := c.ShouldBind(putRequest); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceAdDetails := service.NewAdDetailService()
	serviceCat := service.NewCatService()
	isDisabledCat := 0
	isDisabledAd := 0
	isApprovedAd := 0

	adID, err := manager.SToUint64(sAdID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	ad, err := serviceAds.GetAdByID(adID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, err.Error())
		return

	} else if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	//----------------------------------------------------------
	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserID.Error())
		return
	}

	userRole, ok := c.MustGet("userRole").(string)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserRole.Error())
		return
	}

	if userRole != manager.IsAdmin {
		if ad.UserID != userID {
			c.JSON(http.StatusForbidden, manager.ErrAccessDined.Error())
			return
		}
	}
	//----------------------------------------------------------

	newTitle := strings.TrimSpace(putRequest.Title)
	newCatID := putRequest.CatID
	newPrice := putRequest.Price
	newDescription := strings.TrimRight(putRequest.Description, " ")
	newYoutube := manager.GetYoutubeHash(putRequest.Youtube)
	newLatitude := putRequest.Latitude
	newLongitude := putRequest.Longitude
	newCityName := strings.TrimSpace(putRequest.CityName)
	newCountryName := strings.TrimSpace(putRequest.CountryName)
	newPhoneID := putRequest.PhoneID

	// подменим только то что изменилось
	if ad.Title != newTitle {
		ad.Title = newTitle
	}
	if ad.CatID != newCatID {
		ad.CatID = newCatID
	}
	if ad.Price != newPrice {
		ad.Price = newPrice
	}
	if ad.Description != newDescription {
		ad.Description = newDescription
	}
	if ad.Youtube != newYoutube {
		ad.Youtube = newYoutube
	}
	if ad.Latitude != newLatitude {
		ad.Latitude = newLatitude
	}
	if ad.Longitude != newLongitude {
		ad.Longitude = newLongitude
	}
	if ad.CityName != newCityName {
		ad.CityName = newCityName
	}
	if ad.CountryName != newCountryName {
		ad.CountryName = newCountryName
	}
	if ad.PhoneID != newPhoneID {
		ad.PhoneID = newPhoneID
	}

	if userRole == manager.IsAdmin {
		isDisabledCat = -1
		isDisabledAd = -1
		isApprovedAd = -1

		ad.UserID = putRequest.UserID
		ad.IsDisabled = putRequest.IsDisabled
		ad.IsApproved = putRequest.IsApproved

	} else {
		ad.IsDisabled = false
		ad.IsApproved = false // отправляем на модерацию (принудительно). Если Админ или Модератор, то пришедшее значение, а иначе false.
		ad.IP = c.ClientIP()  // обновим значение
		ad.UserID = userID
	}

	if statusCode, err := checkPhoneNumber(ad); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(statusCode, err.Error())
		return
	}

	catFull, err := serviceCat.GetCatFullByID(ad.CatID, false, isDisabledCat)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// первая проверка объявления (на корректность)
	if err := firstCheck(ad, catFull.Cat); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	// если автогинерация, то подменим. Заголовок обязательно должен быть.
	if catFull.IsAutogenerateTitle {
		updatedTitle, err := updateTitle(catFull, &c.Request.PostForm)
		if err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(400, err.Error())
			return

		} else if updatedTitle == "" {
			c.JSON(400, manager.ErrTitleIsEmpty.Error())
			return
		}

		ad.Title = updatedTitle
	}

	tx := server.Db.Begin()

	if err = serviceAds.Update(ad, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// возьмем текущие фото
	images, err := serviceImages.GetImagesByElIDsAndOpt([]uint64{ad.AdID}, "ad")
	if err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	mFilesAlreadyHas := c.PostFormMap("filesAlreadyHas")

	// если есть разница, то удалим не нужное
	if len(images) != len(mFilesAlreadyHas) {
		restOfImages := make([]*storage.Image, 0) // срез остатков фото

		for _, image := range images {
			var has bool
			for _, requestImageFile := range mFilesAlreadyHas {
				if requestImageFile == image.Filepath {
					restOfImages = append(restOfImages, image)
					has = true
					break
				}
			}

			if !has {
				if err := serviceImages.Delete(image, tx); err != nil {
					tx.Rollback()
					logger.Warning.Println(err.Error())
					c.JSON(500, err.Error())
					return
				}
			}
		}

		images = restOfImages
	}

	if err := workWithPhoto(ad, images, catFull.PropsFull, tx, form, c.SaveUploadedFile); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	adDetailsNew, err := serviceAdDetails.BuildDataFromRequestFormAndCatProps(
		ad.AdID, &c.Request.PostForm, catFull.PropsFull, images)
	if err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	if err := serviceAdDetails.Update(ad.AdID, adDetailsNew, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	tx.Commit()

	adFull, err := serviceAds.GetAdFullByID(ad.AdID, isDisabledAd, isApprovedAd)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	var prefixTest string
	if configs.Cfg.Mode == "debug" {
		prefixTest = ", test"
	}

	emailRequest := emailer.NewEmailRequest(
		configs.Cfg.AdminEmail,
		"Altair (update"+prefixTest+"): "+adFull.Title,
		adFull.Description)
	if ok, err := emailRequest.SendMail(); err != nil {
		logger.Warning.Println(err.Error())

	} else if !ok {
		logger.Warning.Println(manager.ErrEmailNotSend.Error())
	}

	c.JSON(200, adFull)
}

// DeleteAdsAdID - удаление объявления
func DeleteAdsAdID(c *gin.Context) {
	sAdID := c.Param("adID")
	serviceAds := service.NewAdService()

	adID, err := manager.SToUint64(sAdID)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	//----------------------------------------------------------

	userID, ok := c.MustGet("userID").(uint64)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserID.Error())
		return
	}

	userRole, ok := c.MustGet("userRole").(string)
	if !ok {
		c.JSON(400, manager.ErrUndefinedUserRole.Error())
		return
	}

	// тут необходимо проверить на принадлежность данного объявления к текущему пользователю
	if userRole != manager.IsAdmin {
		if ad, err := serviceAds.GetAdByID(adID); err != nil {
			logger.Warning.Println(err.Error())
			c.JSON(500, err.Error())
			return
		} else if ad.UserID != userID {
			c.JSON(400, manager.ErrNotCorrectUserID.Error())
			return
		}
	}

	//----------------------------------------------------------

	tx := server.Db.Begin()
	if err := serviceAds.Delete(adID, tx); err != nil {
		tx.Rollback()
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}
	tx.Commit()

	c.JSON(204, nil)
}

// private -------------------------------------------------------------------------------------------------------------
func workWithPhoto(ad *storage.Ad, curImages []*storage.Image, propsFull []*response.PropFull, tx *gorm.DB,
	form *multipart.Form, fnUpload func(file *multipart.FileHeader, filePath string) error) error {

	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceMediafire := mediafire.NewMediafireService()
	sliceImageSIDs := make([]string, 0)

	// получить текущие фото
	for _, v := range curImages {
		sliceImageSIDs = append(sliceImageSIDs, fmt.Sprint(v.ImgID))
	}

	for _, propFull := range propsFull {
		if propFull.KindPropName != "photo" { // вид св-ва - photo
			continue
		}

		max, err := strconv.Atoi(propFull.Comment)
		if err != nil {
			return err
		}

		max -= len(curImages)

		for key, file := range form.File["files"] {
			// обойдем только определенное число файлов
			if key >= max {
				break
			}

			fileName, err := manager.UploadImage(file, manager.DirImages, fnUpload)
			if err != nil {
				logger.Warning.Println(err.Error())
				continue
			}

			externalFilePath, err := serviceMediafire.UploadSimple(manager.DirImages + "/" + fileName)
			if err != nil {
				logger.Warning.Println(err.Error())
				continue
			} else if externalFilePath == "" {
				logger.Warning.Println(manager.ErrNotFoundExternalFilePath.Error())
				continue
			}

			if err := os.Remove(manager.DirImages + "/" + fileName); err != nil {
				logger.Warning.Println(err.Error())
			}

			image := new(storage.Image)
			image.Filepath = externalFilePath
			image.ElID = ad.AdID
			image.Opt = "ad"

			if err := serviceImages.Create(image, tx); err != nil {
				return err
			}

			sliceImageSIDs = append(sliceImageSIDs, fmt.Sprint(image.ImgID))
		}

		// если есть что добавлять, добавим
		if len(sliceImageSIDs) > 0 {
			// обновим данные у объявления о наличии фото, чтоб удобно потом считать
			ad.HasPhoto = true
			if err := serviceAds.Update(ad, tx); err != nil {
				return err
			}
		}

		break // только один раз (вспомогательная блокировка)
	}

	return nil
}
func firstCheck(ad *storage.Ad, cat *storage.Cat) error {
	serviceCat := service.NewCatService()

	tmpDescription := strings.TrimSpace(ad.Description)
	if tmpDescription == "" {
		return manager.ErrDescIsEmpty
	}

	if isLeaf, err := serviceCat.IsLeaf(ad.CatID); err != nil {
		return err

	} else if !isLeaf {
		return manager.ErrNotCorrectCat
	}

	if !cat.IsAutogenerateTitle && ad.Title == "" {
		return manager.ErrTitleIsRequire
	}

	return nil
}
func updateTitle(catFull *response.СatFull, postForm *url.Values) (string, error) {
	var title string

	switch catFull.CatID {
	// транспорт / автомобили / с пробегом-новый
	case 61, 62:
		var mark string
		model := strings.TrimSpace(postForm.Get("p94"))
		year := strings.TrimSpace(postForm.Get("p95"))

		markID, err := manager.SToUint64(postForm.Get("p72"))
		if err != nil {
			return title, err
		}

		for _, p := range catFull.PropsFull {
			if p.PropID == 72 {
				for _, v2 := range p.Values {
					if v2.ValueID == markID {
						mark = v2.Title
					}
				}
			}
		}

		title = fmt.Sprintf("%s %s, %s", mark, model, year)

	// недвижимость / квартиры / продам / вторичка-новостройка (2-к квартира (Студия), 20 м², 8/9 эт.)
	// Есть "Свободная планировка"
	case 498, 499:
		amountRoomsID, err := manager.SToUint64(postForm.Get("p44"))
		if err != nil {
			return title, err
		}

		var amountRoomsValue string
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 44 {
				for _, v2 := range p.Values {
					if v2.ValueID == amountRoomsID {
						if v2.ValueID == 663 || v2.ValueID == 658 {
							amountRoomsValue = v2.Title

						} else {
							amountRoomsValue = v2.Title + "-к квартира"
						}
					}
				}

			} else if p.PropID == 80 {
				suffix = p.Suffix
			}
		}

		commonArea := strings.TrimSpace(postForm.Get("p80"))
		floor := strings.TrimSpace(postForm.Get("p78"))
		totalFloor := strings.TrimSpace(postForm.Get("p79"))

		// тут есть свободная планировка
		title = fmt.Sprintf("%s, %s %s, %s/%s эт.", amountRoomsValue, commonArea, suffix, floor, totalFloor)

	// недвижимость / квартиры / сдам / на длительный срок-посуточно (2-к квартира (Студия), 20 м², 8/9 эт.)
	case 500, 501:
		amountRoomsID, err := manager.SToUint64(postForm.Get("p46"))
		if err != nil {
			return title, err
		}

		var amountRoomsValue string
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 46 {
				for _, v2 := range p.Values {
					if v2.ValueID == amountRoomsID {
						if v2.ValueID == 676 {
							amountRoomsValue = v2.Title

						} else {
							amountRoomsValue = v2.Title + "-к квартира"
						}
					}
				}

			} else if p.PropID == 80 {
				suffix = p.Suffix
			}
		}

		commonArea := strings.TrimSpace(postForm.Get("p80"))
		floor := strings.TrimSpace(postForm.Get("p78"))
		totalFloor := strings.TrimSpace(postForm.Get("p79"))

		title = fmt.Sprintf("%s, %s %s, %s/%s эт.", amountRoomsValue, commonArea, suffix, floor, totalFloor)

	// недвижимость / квартиры / куплю (Куплю 2-к квартиру (студию и т.д.))
	case 104:
		amountRoomsID, err := manager.SToUint64(postForm.Get("p44"))
		if err != nil {
			return title, err
		}

		var amountRoomsValue string
		for _, p := range catFull.PropsFull {
			if p.PropID == 44 {
				for _, v2 := range p.Values {
					if v2.ValueID == amountRoomsID {
						switch v2.ValueID {
						case 663:
							amountRoomsValue = "студию"
						case 658:
							amountRoomsValue = "свободную планировку"
						default:
							amountRoomsValue = v2.Title + "-к квартиру"
						}
					}
				}
			}
		}

		title = fmt.Sprintf("Куплю %s", amountRoomsValue)

	// недвижимость / квартиры / сниму / на длительный срок-посуточно (Cниму 2-к квартиру (студию))
	case 502, 503:
		amountRoomsID, err := manager.SToUint64(postForm.Get("p46"))
		if err != nil {
			return title, err
		}

		var amountRoomsValue string
		for _, p := range catFull.PropsFull {
			if p.PropID == 46 {
				for _, v2 := range p.Values {
					if v2.ValueID == amountRoomsID {
						if v2.ValueID == 676 {
							amountRoomsValue = "студию"

						} else {
							amountRoomsValue = v2.Title + "-к квартиру"
						}
					}
				}
			}
		}

		title = fmt.Sprintf("Cниму %s", amountRoomsValue)

	// недвижимость / комнаты / продам (Комната 20 м² в 9-к, 8/9 эт.)
	// недвижимость / комнаты / сдам / на длительный срок-посуточно
	case 106, 504, 505:
		roomsID, err := manager.SToUint64(postForm.Get("p49"))
		if err != nil {
			return title, err
		}

		var roomsValue string
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 49 {
				for _, v2 := range p.Values {
					if v2.ValueID == roomsID {
						roomsValue = v2.Title
					}
				}

			} else if p.PropID == 85 {
				suffix = p.Suffix
			}
		}

		areaRoom := strings.TrimSpace(postForm.Get("p85"))
		floor := strings.TrimSpace(postForm.Get("p78"))
		totalFloor := strings.TrimSpace(postForm.Get("p79"))
		title = fmt.Sprintf("Комната %s %s в %s-к, %s/%s эт.", areaRoom, suffix, roomsValue, floor, totalFloor)

	// недвижимость / комнаты / куплю
	case 108:
		title = "Куплю комнату"

	// недвижимость / комнаты / сниму / на длительный срок-посуточно
	case 506, 507:
		title = "Сниму комнату"

	// недвижимость / дома... / продам
	// недвижимость / дома... / сдам / на длительный срок-посуточно
	// (Дом, Дача, Коттедж, Таунхаус) 95 м² на участке 12 сот.
	case 110, 508, 509:
		buildID, err := manager.SToUint64(postForm.Get("p11"))
		if err != nil {
			return title, err
		}

		var buildValue string
		var suffixBuild string
		var suffixEarth string
		for _, p := range catFull.PropsFull {
			switch p.PropID {
			case 11:
				for _, v2 := range p.Values {
					if v2.ValueID == buildID {
						buildValue = v2.Title
					}
				}
			case 86:
				suffixBuild = p.Suffix
			case 87:
				suffixEarth = p.Suffix
			}
		}

		areaBuild := strings.TrimSpace(postForm.Get("p86"))
		areaEarth := strings.TrimSpace(postForm.Get("p87"))
		title = fmt.Sprintf("%s %s %s на участке %s %s", buildValue, areaBuild, suffixBuild, areaEarth, suffixEarth)

	// недвижимость / дома... / куплю. Куплю дом (...)
	case 112:
		buildID, err := manager.SToUint64(postForm.Get("p11"))
		if err != nil {
			return title, err
		}

		var buildValue string
		for _, p := range catFull.PropsFull {
			if p.PropID == 11 {
				for _, v2 := range p.Values {
					if v2.ValueID == buildID {
						if buildID == 93 {
							buildValue = "дачу"
						} else {
							buildValue = strings.ToLower(v2.Title)
						}
					}
				}
			}
		}

		title = fmt.Sprintf("Куплю %s", buildValue)

	// недвижимость / дома... / сниму / на длительный срок-посуточно. Сниму дом (...)
	case 510, 511:
		buildID, err := manager.SToUint64(postForm.Get("p11"))
		if err != nil {
			return title, err
		}

		var buildValue string
		for _, p := range catFull.PropsFull {
			if p.PropID == 11 {
				for _, v2 := range p.Values {
					if v2.ValueID == buildID {
						if buildID == 93 {
							buildValue = "дачу"
						} else {
							buildValue = strings.ToLower(v2.Title)
						}
					}
				}
			}
		}

		title = fmt.Sprintf("Сниму %s", buildValue)

	// недвижимость / земельные участки / продам / поселений (ИЖС)
	// недвижимость / земельные участки / сдам / поселений (ИЖС)
	// Участок 50 сот. (ИЖС)
	case 512, 515:
		area := strings.TrimSpace(postForm.Get("p13"))
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 13 {
				suffix = p.Suffix
			}
		}

		title = fmt.Sprintf("Участок %s %s (ИЖС)", area, suffix)

	// недвижимость / земельные участки / продам / сельхозназначений (СНТ, ДНП)
	// недвижимость / земельные участки / сдам / сельхозназначений (СНТ, ДНП)
	// Участок 50 сот. (СНТ, ДНП)
	case 513, 516:
		area := strings.TrimSpace(postForm.Get("p13"))
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 13 {
				suffix = p.Suffix
			}
		}

		title = fmt.Sprintf("Участок %s %s (СНТ, ДНП)", area, suffix)

	// недвижимость / земельные участки / продам / промназначения
	// недвижимость / земельные участки / сдам / промназначения
	// Участок 50 сот. (промназначения)
	case 514, 517:
		area := strings.TrimSpace(postForm.Get("p13"))
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 13 {
				suffix = p.Suffix
			}
		}

		title = fmt.Sprintf("Участок %s %s (промназначения)", area, suffix)

	// недвижимость / земельные участки / куплю / поселений (ИЖС)
	case 518:
		title = "Куплю участок (ИЖС)"

	// недвижимость / земельные участки / куплю / сельхозназначения (СНТ, ДНП)
	case 519:
		title = "Куплю участок (СНТ, ДНП)"

	// недвижимость / земельные участки / куплю / промназначения
	case 520:
		title = "Куплю участок (промназначения)"

	// недвижимость / земельные участки / сниму / поселений (ИЖС)
	case 521:
		title = "Сниму участок (ИЖС)"

	// недвижимость / земельные участки / сниму / сельхозназначения (СНТ, ДНП)
	case 522:
		title = "Сниму участок (СНТ, ДНП)"

	// недвижимость / земельные участки / сниму / промназначения
	case 523:
		title = "Сниму участок (промназначения)"

	// недвижимость / гаражи и машиноместа / продам-сдам / гараж
	// Гараж, 24 м²
	case 524, 526:
		areaID, err := manager.SToUint64(postForm.Get("p53"))
		if err != nil {
			return title, err
		}

		var areaValue string
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 53 {
				for _, v2 := range p.Values {
					if v2.ValueID == areaID {
						areaValue = v2.Title
					}
				}

				suffix = p.Suffix
			}
		}

		title = fmt.Sprintf("Гараж, %s %s", areaValue, suffix)

	// недвижимость / гаражи и машиноместа / продам-сдам / машиноместо
	// Машиноместо, 24 м²
	case 525, 527:
		areaID, err := manager.SToUint64(postForm.Get("p53"))
		if err != nil {
			return title, err
		}

		var areaValue string
		var suffix string
		for _, p := range catFull.PropsFull {
			if p.PropID == 53 {
				for _, v2 := range p.Values {
					if v2.ValueID == areaID {
						areaValue = v2.Title
					}
				}

				suffix = p.Suffix
			}
		}

		title = fmt.Sprintf("Машиноместо, %s %s", areaValue, suffix)

	// недвижимость / гаражи и машиноместа / куплю / гараж
	case 528:
		title = "Куплю гараж"

	// недвижимость / гаражи и машиноместа / куплю / машиноместо
	case 529:
		title = "Куплю машиноместо"

	// недвижимость / гаражи и машиноместа / сниму / гараж
	case 530:
		title = "Сниму гараж"

	// недвижимость / гаражи и машиноместа / сниму / машиноместо
	case 531:
		title = "Сниму машиноместо"
	}

	// недвижимость / недвижимость за рубежом / ...
	if catFull.CatID >= 561 && catFull.CatID <= 580 {
		countryID, err := manager.SToUint64(postForm.Get("p55"))
		if err != nil {
			return title, err
		}

		var countryValue string
		for _, p := range catFull.PropsFull {
			if p.PropID == 55 {
				for _, v2 := range p.Values {
					if v2.ValueID == countryID {
						countryValue = v2.Title
					}
				}
			}
		}

		if catFull.CatID >= 561 && catFull.CatID <= 570 {
			title = fmt.Sprintf("%s (%s)", catFull.Name, countryValue)
		} else {
			switch catFull.CatID {
			case 571:
				title = fmt.Sprintf("Куплю квартиру (%s)", countryValue)
			case 572:
				title = fmt.Sprintf("Куплю дом (%s)", countryValue)
			case 573:
				title = fmt.Sprintf("Куплю участок (%s)", countryValue)
			case 574:
				title = fmt.Sprintf("Куплю гараж, машиноместо (%s)", countryValue)
			case 575:
				title = fmt.Sprintf("Куплю коммерческую недвижимость (%s)", countryValue)
			case 576:
				title = fmt.Sprintf("Сниму квартиру (%s)", countryValue)
			case 577:
				title = fmt.Sprintf("Сниму дом (%s)", countryValue)
			case 578:
				title = fmt.Sprintf("Сниму участок (%s)", countryValue)
			case 579:
				title = fmt.Sprintf("Сниму гараж, машиноместо (%s)", countryValue)
			case 580:
				title = fmt.Sprintf("Сниму коммерческую недвижимость (%s)", countryValue)
			}
		}
	}

	return title, nil
}
func checkPhoneNumber(ad *storage.Ad) (int, error) {
	servicePhone := service.NewPhoneService()

	phones, err := servicePhone.GetPhonesByUserID(ad.UserID) // id пользователя, либо id польз. что указал Админ
	if err != nil {
		return 500, err

	} else if len(phones) < 1 {
		return 400, manager.ErrPhoneNotFound
	}

	var isSetTruePhoneNumber bool
	for _, v := range phones {
		if v.PhoneID == ad.PhoneID {
			isSetTruePhoneNumber = true
			break
		}
	}
	if !isSetTruePhoneNumber {
		return 400, manager.ErrPhoneSetNotCorrectNumber
	}

	return 200, nil
}
