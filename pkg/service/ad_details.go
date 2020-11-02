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

// NewAdDetailService - фабрика, создает объект деталей объявления
func NewAdDetailService() *AdDetailService {
	return new(AdDetailService)
}

// AdDetailService - структура деталей объявления
type AdDetailService struct{}

// GetByAdID - получить данные относительно id объявления
func (ads AdDetailService) GetByAdID(adID uint64) ([]*storage.AdDetail, error) {
	list := make([]*storage.AdDetail, 0)
	err := server.Db.Where("ad_id = ?", adID).Find(&list).Error

	if err != nil {
		return list, err
	}

	return list, nil
}

// GetByAdIDs - получит данные относительно от нескольких id объявлений
func (ads AdDetailService) GetByAdIDs(adID []uint64) ([]*storage.AdDetail, error) {
	list := make([]*storage.AdDetail, 0)
	err := server.Db.Where("ad_id IN (?)", adID).Find(&list).Error

	if err != nil {
		return list, err
	}

	return list, nil
}

// GetDetailsExtByAdIDs - получить расширенные детали относительно от неск-ких id объявлений
func (ads AdDetailService) GetDetailsExtByAdIDs(adIDs []uint64) ([]*response.AdDetailExt, error) {
	list := make([]*response.AdDetailExt, 0)

	if len(adIDs) < 1 {
		return list, nil
	}

	query := `
		SELECT  AD.ad_id, 
				AD.prop_id, 
				AD.value, 
				P.name AS prop_name,
				(SELECT name FROM kind_props WHERE kind_prop_id = P.kind_prop_id) AS kind_prop_name,
				(SELECT title FROM value_props 
					WHERE	value_id = AD.value AND 
							prop_id = AD.prop_id AND 
							P.kind_prop_id IN (SELECT kind_prop_id FROM kind_props WHERE name="radio" OR name="select")
				) AS value_name
			FROM ad_details AD
				LEFT JOIN props P ON P.prop_id = AD.prop_id
			WHERE AD.ad_id`

	if len(adIDs) == 1 {
		query += " = ?"

	} else {
		query += " IN (?)"
	}

	if err := server.Db.Raw(query, adIDs).Scan(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}

// Create - создание деталей
func (ads AdDetailService) Create(list []*storage.AdDetail, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	for _, adDetail := range list {
		if err := tx.Create(adDetail).Error; err != nil {
			return err // manager.ErrNotCreateNewAdDetail
		}
	}

	return nil
}

// Update - изменение детали
func (ads AdDetailService) Update(adID uint64, list []*storage.AdDetail, tx *gorm.DB) error {
	if err := ads.DeleteAllByAdID(adID, tx); err != nil {
		return err
	}
	if err := ads.Create(list, tx); err != nil {
		return err
	}

	return nil
}

// DeleteAllByAdID - удаление деталей относильено от id объявления
func (ads AdDetailService) DeleteAllByAdID(adID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Where("ad_id = ?", adID).Delete(storage.AdDetail{}).Error

	if err != nil {
		return err
	}

	return nil
}

// BuildDataFromRequestFormAndCatProps - сформировать данные из запроса формы и св-в категорий
func (ads AdDetailService) BuildDataFromRequestFormAndCatProps(adID uint64,
	postForm *url.Values, propsFull []*response.PropFull, images []*storage.Image) ([]*storage.AdDetail, error) {

	adDetails := make([]*storage.AdDetail, 0)

	for _, prop := range propsFull {
		sValue := strings.TrimSpace(postForm.Get(fmt.Sprintf("p%d", prop.PropID))) // именуем по своему (для сжатия данных)
		kind := prop.KindPropName

		if kind == "checkbox" || kind == "radio" || kind == "select" {
			tmpSValue := sValue
			sValue = "" // очистим, необходимо позже вставим (если найдем)

			if tmpSValue != "" {
				iValue, err := manager.SToUint64(tmpSValue)
				if err != nil {
					return adDetails, err
				}

				// посмотрим: есть ли в наличии данное значение
				for _, val := range prop.Values {
					if val.ValueID == iValue {
						sValue = fmt.Sprint(val.ValueID) // берем именно значение, а не всю структуру
						break
					}
				}
			}

		} else if kind == "photo" {
			// если уже есть фото(старые), то установим значение = кол-во фоток.
			if len(images) > 0 {
				sValue = fmt.Sprint(len(images))
			}
		}

		// проверим на обязательное св-во. Если есть уже картинки, то пропустить этот момент
		if prop.IsRequire && sValue == "" {
			return adDetails, fmt.Errorf("поле «%s» обязательно", prop.Title)
		}
		if sValue == "" {
			continue
		}

		adDetail := new(storage.AdDetail)
		adDetail.AdID = adID
		adDetail.PropID = prop.PropID
		adDetail.Value = sValue
		adDetails = append(adDetails, adDetail)
	}

	return adDetails, nil
}
