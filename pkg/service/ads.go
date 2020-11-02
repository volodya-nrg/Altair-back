package service

import (
	"altair/api/response"
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strings"
)

// NewAdService - фабрика, создает объект "объявление"
func NewAdService() *AdService {
	return new(AdService)
}

// AdService - структура объявления
type AdService struct{}

// GetAds - получить все объявления
func (as AdService) GetAds(order string) ([]*storage.Ad, error) {
	ads := make([]*storage.Ad, 0)
	query := server.Db

	if order != "" {
		query = query.Order(order)
	}

	err := query.Find(&ads).Error

	return ads, err
}

// GetAdsFull - получить все объявления (полные)
func (as AdService) GetAdsFull(
	catIDs []uint64,
	checkCountCatIDs bool,
	order string,
	limit int,
	offset uint64,
	isDisabled int,
	isApproved int) ([]*response.AdFull, error) {

	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	var err error

	stm := server.Db.Limit(limit).Offset(int(offset)).Order(order)

	if checkCountCatIDs && len(catIDs) < 1 {
		return adsFull, manager.ErrEmptyListCatIDs
	}

	if isDisabled > 0 {
		stm = stm.Where("is_disabled = 1")

	} else if isDisabled == 0 {
		stm = stm.Where("is_disabled = 0")
	}

	if isApproved > 0 {
		stm = stm.Where("is_approved = 1")

	} else if isApproved == 0 {
		stm = stm.Where("is_approved = 0")
	}

	if len(catIDs) > 0 {
		stm = stm.Where("cat_id IN (?)", catIDs)
	}

	if err := stm.Find(&ads).Error; err != nil {
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

// GetAdByID - получить объявление относильно id
func (as AdService) GetAdByID(adID uint64) (*storage.Ad, error) {
	ad := new(storage.Ad)
	err := server.Db.First(ad, adID).Error

	return ad, err
}

// GetLastAdsByOneCat - получить последние объявления относильно конкретной категории
func (as AdService) GetLastAdsByOneCat(limit int) ([]*response.AdFull, error) {
	var err error
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	query := `
		SELECT *
			FROM ads
			WHERE 
				is_disabled = 0 AND
				is_approved = 1 AND
				cat_id = (SELECT cat_id 
							FROM ads 
							WHERE is_disabled = 0 AND is_approved = 1 
							ORDER BY created_at DESC 
							LIMIT 1)
			ORDER BY created_at DESC
			LIMIT ?`

	if err := server.Db.Raw(query, limit).Scan(&ads).Error; err != nil {
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

// GetAdFullByID - получить полное объявление по его ID
func (as AdService) GetAdFullByID(adID uint64, isDisabled, isApproved int) (*response.AdFull, error) {
	ad := new(storage.Ad)
	adFull := new(response.AdFull)
	stm := server.Db

	if isDisabled > 0 {
		stm = stm.Where("is_disabled = 1")

	} else if isDisabled == 0 {
		stm = stm.Where("is_disabled = 0")
	}

	if isApproved > 0 {
		stm = stm.Where("is_approved = 1")

	} else if isApproved == 0 {
		stm = stm.Where("is_approved = 0")
	}

	if err := stm.First(ad, adID).Error; err != nil {
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

// GetAdsFullBySearchTitle - получить полные объявления по поиску заголовка
func (as AdService) GetAdsFullBySearchTitle(title string, catID, limit, offset uint64, mGetParams url.Values) ([]*response.AdFull, error) {
	serviceProps := NewPropService()
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)

	// map[_:[1585741609119] catId:[53] color:[25] ferma:[asd] q:[title]]
	if catID < 1 {
		return getAdsFullByTitleAndCatID(title, limit, offset, 0)
	}

	// возьмем у данного каталога его св-ва со значениями
	listPPropsFull, err := serviceProps.GetPropsFullByCatID(catID, true)
	if err != nil {
		return adsFull, err
	}

	// посмотрим какие св-ва пришли из вне
	mPropAndValue := make(map[uint64]uint64) // ключ - propId, значение - valueId
	for _, prop := range listPPropsFull {
		if sVal, ok := mGetParams[prop.Name]; ok && len(sVal) > 0 {
			sVal := sVal[0]
			// если это какой либо список, то необходимо проверить на соответствие значений св-ва (их валидность)
			if ok, _ := manager.InArray(prop.KindPropName, manager.TagKindNumber); ok {
				if iVal, err := manager.SToUint64(sVal); err == nil {
					for _, val := range prop.Values {
						if val.ValueID == iVal {
							mPropAndValue[prop.PropID] = iVal
							break
						}
					}
				}
			}
		}
	}

	if len(mPropAndValue) < 1 {
		return getAdsFullByTitleAndCatID(title, limit, offset, catID)
	}

	slicePropValFilter := make([]string, 0)
	for propID, valueID := range mPropAndValue {
		str := fmt.Sprintf("VP.prop_id = %d AND VP.value_id = %d", propID, valueID)
		slicePropValFilter = append(slicePropValFilter, str)
	}

	queryDop := strings.Join(slicePropValFilter, " OR ")
	query := `
		SELECT A.*
			FROM ads A
				LEFT JOIN cats C ON C.cat_id = A.cat_id
				LEFT JOIN cats_props CP ON CP.cat_id = C.cat_id
				LEFT JOIN value_props VP ON VP.prop_id = CP.prop_id
			WHERE 
				A.is_disabled = 0 AND
				A.is_approved = 1 AND
				C.cat_id = ? AND
				A.title LIKE ? AND
				(` + queryDop + `)
			ORDER BY A.created_at DESC
			LIMIT ? 
			OFFSET ?`

	// LIKE "ABC%" = "ABC[ниже]" < KEY < "ABC[выше]". LIKE "%ABC" не может быть оптимизирован для исп-ия индексов
	if err := server.Db.Raw(query, catID, "%"+title+"%", limit, offset).Scan(&ads).Error; err != nil {
		return adsFull, err
	}

	adsFull, err = buildAdsFullFromAds(ads)
	if err != nil {
		return adsFull, err
	}

	return adsFull, nil
}

// GetAdsFullByUserID - получить полные объявления относильно ID пользователя
func (as AdService) GetAdsFullByUserID(userID uint64, order string, limit, offset uint64) ([]*response.AdFull, error) {
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	var err error

	query := server.Db.Limit(int(limit)).Offset(int(offset)).Order(order).Where("user_id = ?", userID)

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

// GetAdFullByUserIDAndCatID - получить полное объявление относительно ID пользователя и ID категории
func (as AdService) GetAdFullByUserIDAndCatID(userID, adID uint64) (*response.AdFull, error) {
	ad := new(storage.Ad)
	adFull := new(response.AdFull)

	if err := server.Db.Where("user_id = ? AND ad_id = ?", userID, adID).First(ad).Error; err != nil {
		return adFull, err
	}

	// доп. проверка
	if ad.AdID < 1 {
		return adFull, manager.ErrNotFoundAd
	}

	adsFull, err := buildAdsFullFromAds([]*storage.Ad{ad})
	if err != nil {
		return adFull, err
	}

	adFull = adsFull[0]

	return adFull, nil
}

// GetAdsByUserID - получить объявления относительно ID пользователя
func (as AdService) GetAdsByUserID(userID uint64) ([]*storage.Ad, error) {
	ads := make([]*storage.Ad, 0)
	err := server.Db.Where("user_id = ?", userID).Find(&ads).Error

	return ads, err
}

// Create - создать объявление
func (as AdService) Create(ad *storage.Ad, tx *gorm.DB) error {
	uniqStr := manager.RandStringRunes(5)
	ad.Slug = fmt.Sprintf("%s_%s", manager.TranslitRuToEn(ad.Title), uniqStr)

	if tx == nil {
		tx = server.Db
	}
	if err := tx.Create(ad).Error; err != nil {
		return err
	}
	if err := as.Update(ad, tx); err != nil {
		return err
	}

	return nil
}

// Update - обновить обновление
func (as AdService) Update(ad *storage.Ad, tx *gorm.DB) error {
	ad.Slug = fmt.Sprintf("%s_%d", manager.TranslitRuToEn(ad.Title), ad.AdID)

	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(ad).Error

	return err
}

// UpdateByPhoneID - обновить объявление относительно ID телефона
func (as AdService) UpdateByPhoneID(phoneIDOld, phoneIDNew, userID uint64, tx *gorm.DB) error {
	ads := make([]*storage.Ad, 0)
	adIDs := make([]uint64, 0)
	var query string

	if tx == nil {
		tx = server.Db
	}

	err := tx.Where("phone_id = ? AND user_id = ?", phoneIDOld, userID).Find(&ads).Error
	if err != nil {
		return err
	}
	if len(ads) < 1 {
		return nil
	}

	for _, v := range ads {
		adIDs = append(adIDs, v.AdID)
	}

	// заменить id номера телефона на валидный (пришедший)
	if phoneIDNew > 0 {
		query = `UPDATE ads SET phone_id = ? WHERE ad_id IN (?)`
	} else {
		query = `UPDATE ads SET phone_id = ?, is_disabled = 1, is_approved = 0 WHERE ad_id IN (?)`
	}

	if err := tx.Exec(query, phoneIDNew, adIDs).Error; err != nil {
		return err
	}

	return nil
}

// Delete - удалить объявление
func (as AdService) Delete(adID uint64, tx *gorm.DB) error {
	serviceImages := NewImageService()
	serviceAdDetail := NewAdDetailService()

	if tx == nil {
		tx = server.Db
	}

	if err := tx.Delete(storage.Ad{}, "ad_id = ?", adID).Error; err != nil {
		return err
	}

	images, err := serviceImages.GetImagesByElIDsAndOpt([]uint64{adID}, "ad")
	if err != nil {
		return err
	}

	for _, v := range images {
		if err := serviceImages.Delete(v, tx); err != nil {
			return err
		}
	}

	if err := serviceAdDetail.DeleteAllByAdID(adID, tx); err != nil {
		return err
	}

	return nil
}

// DeleteAllByUserID - удалить все объявления относительно ID пользователя
func (as AdService) DeleteAllByUserID(userID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	ads, err := as.GetAdsByUserID(userID)
	if err != nil {
		return err
	}

	for _, v := range ads {
		if err := as.Delete(v.AdID, tx); err != nil {
			return err
		}
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func getAdsFullByTitleAndCatID(title string, limitSrc, offsetSrc, catID uint64) ([]*response.AdFull, error) {
	ads := make([]*storage.Ad, 0)
	adsFull := make([]*response.AdFull, 0)
	var err error
	var limit uint64 = 10
	var offset uint64 = 0

	if limitSrc > 0 && limitSrc < 10 {
		limit = limitSrc
	}
	if offsetSrc > 0 {
		offset = offsetSrc
	}

	stm := server.Db.Limit(int(limit)).
		Where("is_disabled = 0 AND is_approved = 1 AND title LIKE ?", "%"+title+"%")

	if offset > 0 {
		stm = stm.Offset(int(offset))
	}

	if catID > 0 {
		stm = stm.Where("cat_id = ?", catID)
	}

	if err := stm.Find(&ads).Error; err != nil {
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
	adIDs := make([]uint64, 0)

	if len(ads) < 1 {
		return adsFull, nil
	}

	for _, ad := range ads {
		adIDs = append(adIDs, ad.AdID)
		adFull := new(response.AdFull)
		adFull.Ad = ad

		// с массивами надо сделать отдельную ф-ии фабрику, где по умолчанию в массивах
		// будет стоять пустой массив
		adFull.Images = make([]*storage.Image, 0)
		adFull.DetailsExt = make([]*response.AdDetailExt, 0)

		adsFull = append(adsFull, adFull)
	}

	images, err := serviceImages.GetImagesByElIDsAndOpt(adIDs, "ad")
	if err != nil {
		return adsFull, err
	}

	adDetailsExt, err := serviceAdDetails.GetDetailsExtByAdIDs(adIDs)
	if err != nil {
		return adsFull, err
	}

	for _, adFull := range adsFull {
		for _, image := range images {
			if image.ElID == adFull.AdID {
				adFull.Images = append(adFull.Images, image)
			}
		}
		for _, detailExt := range adDetailsExt {
			if detailExt.AdID == adFull.AdID {
				adFull.DetailsExt = append(adFull.DetailsExt, detailExt)
			}
		}
	}

	return adsFull, nil
}
