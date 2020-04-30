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

func (as AdService) GetAds(order string) ([]*storage.Ad, error) {
	ads := make([]*storage.Ad, 0)
	query := server.Db.Debug()

	if order != "" {
		query = query.Order(order)
	}

	err := query.Find(&ads).Error

	return ads, err
}
func (as AdService) GetAdsFull(catIds []uint64, checkCountCatIds bool, order string, limit uint64, offset uint64) ([]*response.AdFull, error) {
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	var err error

	query := server.Db.Debug().Limit(limit).Offset(offset).Order(order)

	if checkCountCatIds && len(catIds) < 1 {
		return adsFull, errEmptyListCatIds
	}

	if len(catIds) > 0 {
		query = query.Where("cat_id IN (?)", catIds)
	}

	if err := query.Find(&ads).Error; err != nil {
		return adsFull, err
	}

	if len(ads) < 1 {
		return adsFull, nil
	}

	adsFull, err = buildAdsFullFromAds(ads)
	if err != nil {
		return adsFull, err
	}

	return adsFull, nil
}
func (as AdService) GetAdById(adId uint64) (*storage.Ad, error) {
	ad := new(storage.Ad)
	err := server.Db.Debug().First(ad, adId).Error

	return ad, err
}
func (as AdService) GetLastAdsByOneCat(limit uint64) ([]*response.AdFull, error) {
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	var err error
	query := `
		SELECT * 
			FROM ads
			WHERE cat_id = (SELECT cat_id FROM ads ORDER BY created_at DESC LIMIT 1)
			LIMIT ?`

	if err := server.Db.Debug().Raw(query, limit).Scan(&ads).Error; err != nil {
		return adsFull, err
	}

	if len(ads) < 1 {
		return adsFull, nil
	}

	adsFull, err = buildAdsFullFromAds(ads)
	if err != nil {
		return adsFull, err
	}

	return adsFull, err
}
func (as AdService) GetAdFullById(adId uint64) (*response.AdFull, error) {
	ad := new(storage.Ad)
	adFull := new(response.AdFull)

	if err := server.Db.Debug().First(ad, adId).Error; err != nil {
		return adFull, err
	}

	adsFull, err := buildAdsFullFromAds([]*storage.Ad{ad})
	if err != nil {
		return adFull, err
	}
	if len(adsFull) < 1 {
		return adFull, nil
	}

	return adsFull[0], nil
}
func (as AdService) GetAdsFullBySearchTitle(title string, catId uint64, limit uint64, offset uint64, mGetParams url.Values) ([]*response.AdFull, error) {
	serviceProps := NewPropService()
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)

	// map[_:[1585741609119] catId:[53] color:[25] ferma:[asd] q:[title]]
	if catId < 1 {
		return getAdsFullByTitleAndCatId(title, limit, offset, 0)
	}

	// возьмем у данного каталога его св-ва со значениями
	listPPropsFull, err := serviceProps.GetPropsFullByCatId(catId, true)
	if err != nil {
		return adsFull, err
	}

	// посмотрим какие св-ва пришли из вне
	mPropAndValue := make(map[uint64]uint64, 0) // ключ - propId, значение - valueId
	for _, prop := range listPPropsFull {
		if sVal, ok := mGetParams[prop.Name]; ok && len(sVal) > 0 {
			sVal := sVal[0]
			// если это какой либо список, то необходимо проверить на соответствие значений св-ва (их валидность)
			if prop.KindPropName == "radio" || prop.KindPropName == "select" {
				if iVal, err := strconv.ParseUint(sVal, 10, 64); err == nil {
					for _, val := range prop.Values {
						if val.ValueId == iVal {
							mPropAndValue[prop.PropId] = iVal
							break
						}
					}
				}
			}
		}
	}

	if len(mPropAndValue) < 1 {
		return getAdsFullByTitleAndCatId(title, limit, offset, catId)
	}

	slicePropValFilter := make([]string, 0)
	for propId, valueId := range mPropAndValue {
		str := fmt.Sprintf("VP.prop_id = %d AND VP.value_id = %d", propId, valueId)
		slicePropValFilter = append(slicePropValFilter, str)
	}

	queryDop := strings.Join(slicePropValFilter[:], " OR ")
	query := `
		SELECT A.*
			FROM ads A
				LEFT JOIN cats C ON C.cat_id = A.cat_id
				LEFT JOIN cats_props CP ON CP.cat_id = C.cat_id
				LEFT JOIN value_props VP ON VP.prop_id = CP.prop_id
			WHERE 
				C.cat_id = ? AND
				A.title LIKE ? AND
				(` + queryDop + `)
			ORDER BY A.created_at DESC
			LIMIT 100`
	// LIKE "ABC%" = "ABC[ниже]" < KEY < "ABC[выше]". LIKE "%ABC" не может быть оптимизирован для исп-ия индексов
	if err := server.Db.Debug().Raw(query, catId, "%"+title+"%").Scan(&ads).Error; err != nil {
		return adsFull, err
	}

	adsFull, err = buildAdsFullFromAds(ads)
	if err != nil {
		return adsFull, err
	}

	return adsFull, nil
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
func (as AdService) Delete(adId uint64, tx *gorm.DB) error {
	serviceImages := NewImageService()
	serviceAdDetail := NewAdDetailService()

	if tx == nil {
		tx = server.Db.Debug()
	}

	if err := tx.Delete(storage.Ad{}, "ad_id = ?", adId).Error; err != nil {
		return err
	}

	images, err := serviceImages.GetImagesByElIdsAndOpt([]uint64{adId}, "ad")
	if err != nil {
		return err
	}

	for _, v := range images {
		if err := serviceImages.Delete(v, tx); err != nil {
			return err
		}
	}

	if err := serviceAdDetail.DeleteAllByAdId(adId, tx); err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func getAdsFullByTitleAndCatId(title string, limitSrc uint64, offsetSrc uint64, catId uint64) ([]*response.AdFull, error) {
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	query := server.Db.Debug().Where("title LIKE ?", "%"+title+"%")
	var err error
	var limit uint64 = 10
	var offset uint64 = 0

	if limitSrc > 0 && limitSrc < 10 {
		limit = limitSrc
	}
	if offsetSrc > 0 {
		offset = offsetSrc
	}

	query = query.Limit(limit)

	if offset > 0 {
		query = query.Offset(offset)
	}

	if catId > 0 {
		query = query.Where("cat_id = ?", catId)
	}

	if err := query.Find(&ads).Error; err != nil {
		return adsFull, err
	}

	adsFull, err = buildAdsFullFromAds(ads)
	if err != nil {
		return adsFull, err
	}

	return adsFull, nil
}

func buildAdsFullFromAds(ads []*storage.Ad) ([]*response.AdFull, error) {
	serviceImages := NewImageService()
	serviceAdDetails := NewAdDetailService()
	adsFull := make([]*response.AdFull, 0)
	adIds := make([]uint64, 0)
	catIds := make([]uint64, 0)

	if len(ads) < 1 {
		return adsFull, nil
	}

	for _, ad := range ads {
		adIds = append(adIds, ad.AdId)
		adFull := new(response.AdFull)
		adFull.Ad = ad
		adsFull = append(adsFull, adFull)
		catIds = append(catIds, ad.CatId)
	}

	images, err := serviceImages.GetImagesByElIdsAndOpt(adIds, "ad")
	if err != nil {
		return adsFull, err
	}

	adDetailsExt, err := serviceAdDetails.GetDetailsExtByAdIds(adIds)
	if err != nil {
		return adsFull, err
	}

	for _, adFull := range adsFull {
		for _, image := range images {
			if image.ElId == adFull.AdId {
				adFull.Images = append(adFull.Images, image)
			}
		}
		for _, detailExt := range adDetailsExt {
			if detailExt.AdId == adFull.AdId {
				adFull.DetailsExt = append(adFull.DetailsExt, detailExt)
			}
		}
	}

	return adsFull, nil
}
