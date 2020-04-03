package service

import "C"
import (
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"strconv"
	"strings"
)

func NewAdService() *AdService {
	return new(AdService)
}

type AdService struct{}

func (as AdService) GetAds(isOrderDesc bool) ([]*storage.Ad, error) {
	pAds := make([]*storage.Ad, 0)
	order := "asc"

	if isOrderDesc {
		order = "desc"
	}

	err := server.Db.Debug().Order("created_at " + order).Find(&pAds).Error

	return pAds, err
}
func (as AdService) GetAdsFull(catIds []uint64, checkCountCatIds bool, isOrderDesc bool, serviceImages *ImageService, serviceAdDetails *AdDetailService) ([]*response.AdFull, error) {
	pAds := make([]*storage.Ad, 0)
	pAdsFull := make([]*response.AdFull, 0)
	var err error
	order := "asc"

	if isOrderDesc {
		order = "desc"
	}

	query := server.Db.Debug().Order("created_at " + order)

	if checkCountCatIds && len(catIds) < 1 {
		return pAdsFull, errEmptyListCatIds
	}

	if len(catIds) > 0 {
		query = query.Where("cat_id IN (?)", catIds)
	}

	if err := query.Find(&pAds).Error; err != nil {
		return pAdsFull, err
	}

	if len(pAds) < 1 {
		return pAdsFull, nil
	}

	pAdsFull, err = as.buildAdsFullFromAds(pAds, serviceImages, serviceAdDetails)
	if err != nil {
		return pAdsFull, err
	}

	return pAdsFull, nil
}
func (as AdService) GetAdById(adId uint64) (*storage.Ad, error) {
	pAd := new(storage.Ad)
	err := server.Db.Debug().First(pAd, adId).Error

	return pAd, err
}
func (as AdService) GetAdFullById(adId uint64, serviceImages *ImageService, serviceAdDetails *AdDetailService) (*response.AdFull, error) {
	pAd := new(storage.Ad)
	pAdFull := new(response.AdFull)

	if err := server.Db.Debug().First(pAd, adId).Error; err != nil {
		return pAdFull, err
	}

	pImages, err := serviceImages.GetImagesByElIdsAndOpt([]uint64{pAd.AdId}, "ad")
	if err != nil {
		return pAdFull, err
	}

	pAdDetailsExt, err := serviceAdDetails.GetDetailsExtByAdIds([]uint64{adId})
	if err != nil {
		return pAdFull, err
	}

	pAdFull.Ad = pAd
	pAdFull.Images = pImages
	pAdFull.Details = pAdDetailsExt

	return pAdFull, nil
}
func (as AdService) GetAdFullByIds(adIds []uint64, isOrderDesc bool, serviceImages *ImageService, serviceAdDetails *AdDetailService) ([]*response.AdFull, error) {
	pAds := make([]*storage.Ad, 0)
	pAdFulls := make([]*response.AdFull, 0)
	order := "asc"

	if len(adIds) < 1 {
		return pAdFulls, nil
	}

	if isOrderDesc {
		order = "desc"
	}

	if err := server.Db.Debug().Order("created_at "+order).Find(&pAds, adIds).Error; err != nil {
		return pAdFulls, err
	}

	pImages, err := serviceImages.GetImagesByElIdsAndOpt(adIds, "ad")
	if err != nil {
		return pAdFulls, err
	}

	pAdDetailsExt, err := serviceAdDetails.GetDetailsExtByAdIds(adIds)
	if err != nil {
		return pAdFulls, err
	}

	// теперь три эти составляющие (ad, images, details) необходимо объединить
	for _, ad := range pAds {
		pAdFull := new(response.AdFull)
		pAdFull.Ad = ad

		for _, image := range pImages {
			if image.ElId == ad.AdId {
				pAdFull.Images = append(pAdFull.Images, image)
			}
		}
		for _, detailsExt := range pAdDetailsExt {
			if detailsExt.AdId == ad.AdId {
				pAdFull.Details = append(pAdFull.Details, detailsExt)
			}
		}

		pAdFulls = append(pAdFulls, pAdFull)
	}

	return pAdFulls, nil
}
func (as AdService) GetAdsFullBySearchTitle(title string, catId uint64, mGetParams url.Values, serviceImages *ImageService, serviceAdDetails *AdDetailService, serviceProperties *PropertyService, valuesPropertyService *ValuesPropertyService) ([]*response.AdFull, error) {
	pAdsFull := make([]*response.AdFull, 0)

	// map[_:[1585741609119] catId:[53] color:[25] ferma:[asd] q:[title]]
	if catId < 1 {
		return as.getAdsFullByTitleAndCatId(title, 0, serviceImages, serviceAdDetails)
	}

	// возьмем у данного каталога его св-ва со значениями
	listPPropertiesFull, err := serviceProperties.GetPropertiesFullByCatId(catId, true, valuesPropertyService)
	if err != nil {
		return pAdsFull, err
	}

	// посмотрим какие св-ва пришли из вне
	mPropertyAndValue := make(map[uint64]uint64, 0) // ключ - propertyId, значение - valueId
	for _, prop := range listPPropertiesFull {
		if sVal, ok := mGetParams[prop.Name]; ok && len(sVal) > 0 {
			sVal := sVal[0]
			// если это какой либо список, то необходимо проверить на соответствие значений св-ва (их валидность)
			if prop.KindPropertyName == "radio" || prop.KindPropertyName == "select" {
				if iVal, err := strconv.ParseUint(sVal, 10, 64); err == nil {
					for _, val := range prop.Values {
						if val.ValueId == iVal {
							mPropertyAndValue[prop.PropertyId] = iVal
							break
						}
					}
				}
			}
		}
	}

	if len(mPropertyAndValue) < 1 {
		return as.getAdsFullByTitleAndCatId(title, catId, serviceImages, serviceAdDetails)
	}

	slicePropValFilter := make([]string, 0)
	for propertyId, valueId := range mPropertyAndValue {
		str := fmt.Sprintf("VP.property_id = %d AND VP.value_id = %d", propertyId, valueId)
		slicePropValFilter = append(slicePropValFilter, str)
	}

	pAds := make([]*storage.Ad, 0)
	adIds := make([]uint64, 0)
	queryDop := strings.Join(slicePropValFilter[:], " OR ")
	query := `
		SELECT A.*
			FROM ads A
				LEFT JOIN cats C ON C.cat_id = A.cat_id
				LEFT JOIN cats_properties CP ON CP.cat_id = C.cat_id
				LEFT JOIN value_properties VP ON VP.property_id = CP.property_id
			WHERE 
				C.cat_id = ? AND
				A.title LIKE ? AND
				(` + queryDop + `)
			ORDER BY A.created_at DESC
			LIMIT 100`

	if err := server.Db.Debug().Raw(query, catId, "%"+title+"%").Scan(&pAds).Error; err != nil {
		return pAdsFull, err
	}

	for _, pAd := range pAds {
		adIds = append(adIds, pAd.AdId)
	}

	return as.GetAdFullByIds(adIds, true, serviceImages, serviceAdDetails)
}
func (as AdService) Create(ad *storage.Ad, tx *gorm.DB) error {
	uniqStr := helpers.RandStringRunes(5)
	ad.Slug = fmt.Sprintf("%s_%s", helpers.TranslitRuToEn(ad.Title), uniqStr)

	if tx == nil {
		tx = server.Db.Debug()
	}
	if !server.Db.Debug().NewRecord(ad) {
		return errOnNewRecordNewAd
	}
	if err := tx.Create(ad).Error; err != nil {
		return errNotCreateNewAd
	}

	err := as.Update(ad, tx)
	if err != nil {
		return err
	}

	return nil
}
func (as AdService) Update(ad *storage.Ad, tx *gorm.DB) error {
	ad.Slug = fmt.Sprintf("%s_%d", helpers.TranslitRuToEn(ad.Title), ad.AdId)

	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(ad).Error

	return err
}
func (as AdService) Delete(adId uint64, tx *gorm.DB, serviceImages *ImageService) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	if err := tx.Where("ad_id = ?", adId).Delete(storage.Ad{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	pImages, err := serviceImages.GetImagesByElIdsAndOpt([]uint64{adId}, "ad")
	if err != nil {
		return err
	}

	if err := serviceImages.DeleteAll(pImages); err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func (as AdService) buildAdsFullFromAds(pAds []*storage.Ad, serviceImages *ImageService, serviceAdDetails *AdDetailService) ([]*response.AdFull, error) {
	pAdsFull := make([]*response.AdFull, 0)
	adIds := make([]uint64, 0)

	if len(pAds) < 1 {
		return pAdsFull, nil
	}

	for _, ad := range pAds {
		adIds = append(adIds, ad.AdId)
		pAdFull := new(response.AdFull)
		pAdFull.Ad = ad
		pAdsFull = append(pAdsFull, pAdFull)
	}

	pImages, err := serviceImages.GetImagesByElIdsAndOpt(adIds, "ad")
	if err != nil {
		return pAdsFull, err
	}

	pAdDetails, err := serviceAdDetails.GetDetailsExtByAdIds(adIds)
	if err != nil {
		return pAdsFull, err
	}

	for _, ad := range pAdsFull {
		for _, image := range pImages {
			if image.ElId == ad.AdId {
				ad.Images = append(ad.Images, image)
			}
		}
		for _, detail := range pAdDetails {
			if detail.AdId == ad.AdId {
				ad.Details = append(ad.Details, detail)
			}
		}
	}

	return pAdsFull, nil
}
func (as AdService) getAdsFullByTitleAndCatId(title string, catId uint64, serviceImages *ImageService, serviceAdDetails *AdDetailService) ([]*response.AdFull, error) {
	pAds := make([]*storage.Ad, 0)
	pAdsFull := make([]*response.AdFull, 0)
	adIds := make([]uint64, 0)
	query := server.Db.Debug().Where("title LIKE ?", "%"+title+"%")

	if catId > 0 {
		query = query.Where("cat_id = ?", catId)
	}

	if err := query.Find(&pAds).Error; err != nil {
		return pAdsFull, err
	}

	for _, ad := range pAds {
		adIds = append(adIds, ad.AdId)
	}

	return as.GetAdFullByIds(adIds, true, serviceImages, serviceAdDetails)
}
