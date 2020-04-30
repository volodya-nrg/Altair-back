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
func (ads AdDetailService) GetByAdIds(adId []uint64) ([]*storage.AdDetail, error) {
	list := make([]*storage.AdDetail, 0)
	err := server.Db.Debug().Where("ad_id IN (?)", adId).Find(&list).Error

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
		SELECT  AD.ad_id, 
				AD.prop_id, 
				AD.value, 
				P.name AS prop_name,
				(SELECT name FROM kind_props WHERE kind_prop_id = P.kind_prop_id) AS kind_prop_name,
				(SELECT title FROM value_props 
					WHERE P.kind_prop_id IN (
						SELECT kind_prop_id FROM kind_props WHERE name="radio" OR name="select" 
					) AND value_id = AD.value AND prop_id = AD.prop_id) AS value_name
			FROM ad_details AD
				LEFT JOIN props P ON P.prop_id = AD.prop_id
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
		if !tx.NewRecord(adDetail) {
			return errOnNewRecordNewAdDetail
		}
		if err := tx.Create(adDetail).Error; err != nil {
			return errNotCreateNewAdDetail
		}
	}

	return nil
}
func (ads AdDetailService) Update(adId uint64, list []*storage.AdDetail, tx *gorm.DB) error {
	if err := ads.DeleteAllByAdId(adId, tx); err != nil {
		return err
	}

	if err := ads.Create(list, tx); err != nil {
		return err
	}

	return nil
}
func (ads AdDetailService) DeleteAllByAdId(adId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := server.Db.Debug().Where("ad_id = ?", adId).Delete(storage.AdDetail{}).Error

	if err != nil {
		return err
	}

	return nil
}
func (ads AdDetailService) BuildDataFromRequestFormAndCatProps(
	adId uint64, postForm *url.Values, propsFull []*response.PropFull) ([]*storage.AdDetail, error) {
	adDetails := make([]*storage.AdDetail, 0)

	for _, prop := range propsFull {
		sValue := strings.TrimSpace(postForm.Get(prop.Name))
		kind := prop.KindPropName

		// проверим на обязательное св-во
		if prop.IsRequire && sValue == "" {
			return adDetails, fmt.Errorf("prop (%s) is require", kind)
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
				return adDetails, fmt.Errorf("not found valueId(%d) on prop(%s)", iValue, kind)
			}
		}

		adDetail := new(storage.AdDetail)
		adDetail.AdId = adId
		adDetail.PropId = prop.PropId
		adDetail.Value = sValue

		adDetails = append(adDetails, adDetail)
	}

	return adDetails, nil
}
