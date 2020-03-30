package service

import (
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"fmt"
)

func NewAdService() *AdService {
	return new(AdService)
}

type AdService struct{}

func (as AdService) GetAds() ([]*storage.Ad, error) {
	pAds := make([]*storage.Ad, 0)
	err := server.Db.Debug().Order("created_at", true).Find(&pAds).Error

	return pAds, err
}
func (as AdService) GetAdsFull(catIds []uint64, checkCountCatIds bool) ([]*response.AdFull, error) {
	pAds := make([]*storage.Ad, 0)
	pAdsFull := make([]*response.AdFull, 0)
	var err error
	query := server.Db.Debug().Order("created_at", true)

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

	pAdsFull, err = as.createAdsFullFromAds(pAds)
	if err != nil {
		return pAdsFull, err
	}

	return pAdsFull, nil
}
func (as AdService) GetAdById(adId uint64) (*storage.Ad, error) {
	pAd := new(storage.Ad)
	err := server.Db.Debug().First(pAd, adId).Error // проверяется в контроллере

	return pAd, err
}
func (as AdService) GetAdFullById(adId uint64) (*response.AdFull, error) {
	serviceImages := NewImageService()
	serviceAdDetails := NewAdDetailService()
	pAd := new(storage.Ad)
	pAdFull := new(response.AdFull)

	if err := server.Db.Debug().First(pAd, adId).Error; err != nil {
		return pAdFull, err
	}

	pImages, err := serviceImages.GetImages(pAd.AdId, "ad")
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
func (as AdService) Create(ad *storage.Ad) error {
	uniqStr := helpers.RandStringRunes(5)
	ad.Slug = fmt.Sprintf("%s_%s", helpers.TranslitRuToEn(ad.Title), uniqStr)

	if !server.Db.Debug().NewRecord(ad) {
		return errOnNewRecordNewAd
	}
	if err := server.Db.Debug().Create(ad).Error; err != nil {
		return errNotCreateNewAd
	}

	err := as.Update(ad)
	if err != nil {
		return err
	}

	return nil
}
func (as AdService) Update(ad *storage.Ad) error {
	ad.Slug = fmt.Sprintf("%s_%d", helpers.TranslitRuToEn(ad.Title), ad.AdId)

	return server.Db.Debug().Save(ad).Error
}
func (as AdService) Delete(adId uint64) error {
	serviceImages := NewImageService()
	ad := storage.Ad{
		AdId: adId,
	}

	if err := server.Db.Debug().Delete(ad).Error; err != nil {
		return err
	}

	pImages, err := serviceImages.GetImages(adId, "ad")
	if err != nil {
		return err
	}

	if err := serviceImages.DeleteAll(pImages); err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func (as AdService) createAdsFullFromAds(pAds []*storage.Ad) ([]*response.AdFull, error) {
	serviceImages := NewImageService()
	serviceAdDetails := NewAdDetailService()
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
