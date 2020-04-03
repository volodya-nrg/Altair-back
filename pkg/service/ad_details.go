package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"strconv"
	"strings"
)

func NewAdDetailService() *AdDetailService {
	return new(AdDetailService)
}

type AdDetailService struct{}

func (ads AdDetailService) GetByAdId(adId uint64) ([]*storage.AdDetail, error) {
	list := make([]*storage.AdDetail, 0)
	err := server.Db.Debug().Where("ad_id = ?", adId).Find(&list).Error

	if err != nil {
		return list, err
	}

	return list, nil
}
func (ads AdDetailService) GetDetailsExtByAdIds(adIds []uint64) ([]*response.AdDetailExt, error) {
	list := make([]*response.AdDetailExt, 0)

	if len(adIds) < 1 {
		return list, nil
	}

	query := `
		SELECT AD.ad_id, AD.property_id, AD.value, P.name AS property_name,
			(SELECT name FROM kind_properties WHERE kind_property_id = P.kind_property_id) AS kind_property_name
			FROM ad_details AD
				LEFT JOIN properties P ON P.property_id = AD.property_id
			WHERE AD.ad_id`

	if len(adIds) == 1 {
		query += " = ?"

	} else {
		query += " IN (?)"
	}

	if err := server.Db.Debug().Raw(query, adIds).Scan(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}
func (ads AdDetailService) Create(list []*storage.AdDetail, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	for _, adDetail := range list {
		if !server.Db.Debug().NewRecord(adDetail) {
			return errOnNewRecordNewAdDetail
		}
		if err := tx.Create(adDetail).Error; err != nil {
			return errNotCreateNewAdDetail
		}
	}

	return nil
}
func (ads AdDetailService) Update(adId uint64, list []*storage.AdDetail, tx *gorm.DB) error {
	if err := ads.DeleteByAdId(adId, tx); err != nil {
		return err
	}

	if err := ads.Create(list, tx); err != nil {
		return err
	}

	return nil
}
func (ads AdDetailService) DeleteByAdId(adId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := server.Db.Debug().Where("ad_id = ?", adId).Delete(storage.AdDetail{}).Error

	if err != nil {
		return err
	}

	return nil
}
func (ads AdDetailService) BuildDataFromRequestFormAndCatProps(adId uint64, postForm url.Values, propsFull []*response.PropertyFull) ([]*storage.AdDetail, error) {
	adDetails := make([]*storage.AdDetail, 0)

	for _, prop := range propsFull {
		sValue := strings.TrimSpace(postForm.Get(prop.Name))
		kind := prop.KindPropertyName

		// проверим на обязательное св-во
		if prop.IsRequire && sValue == "" {
			return adDetails, fmt.Errorf("property (%s) is require", kind)
		}
		if sValue == "" {
			continue
		}

		if kind == "checkbox" || kind == "radio" || kind == "select" {
			iValue, err := strconv.ParseUint(sValue, 10, 64)
			if err != nil {
				return adDetails, err
			}

			// посмотрим: есть ли в наличии данное значение
			var has bool
			for _, val := range prop.Values {
				if val.ValueId == iValue {
					has = true
				}
			}

			if !has {
				return adDetails, fmt.Errorf("not found valueId(%d) on property(%s)", iValue, kind)
			}
		}

		pAdDetail := new(storage.AdDetail)
		pAdDetail.AdId = adId
		pAdDetail.PropertyId = prop.PropertyId
		pAdDetail.Value = sValue

		adDetails = append(adDetails, pAdDetail)
	}

	return adDetails, nil
}
