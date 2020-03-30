package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func NewPropertyService() *PropertyService {
	prop := new(PropertyService)
	prop.tblFields = "P.property_id, P.title, P.kind_property_id, P.name, P.is_can_as_filter, P.max_int"

	return prop
}

type PropertyService struct {
	tblFields string
}

func (ps PropertyService) GetProperties(isOrderDesc bool) ([]*storage.Property, error) {
	list := make([]*storage.Property, 0)
	err := server.Db.Debug().Order("property_id", isOrderDesc).Find(&list).Error

	return list, err
}
func (ps PropertyService) GetPropertiesFull() ([]*response.PropertyFull, error) {
	// тут в запросе спец-но не полные данные, т.к. категории нет. Только добавляем kindPropertyName
	list := make([]*response.PropertyFull, 0)
	query := `
		SELECT ` + ps.tblFields + `, KP.name AS kind_property_name
			FROM properties P
				LEFT JOIN kind_properties KP ON KP.kind_property_id = P.kind_property_id
			ORDER BY P.property_id ASC`
	err := server.Db.Debug().Raw(query).Scan(&list).Error

	return list, err
}
func (ps PropertyService) GetPropertiesFullByCatId(catId uint64) ([]*response.PropertyFull, error) {
	valuesPropertyService := NewValuesPropertyService()
	list := make([]*response.PropertyFull, 0)
	query := `
		SELECT  ` + ps.tblFields + `, 
				KP.name AS kind_property_name, 
				CP.pos AS property_pos,
				CP.is_require AS property_is_require
			FROM properties P
				LEFT JOIN kind_properties KP ON KP.kind_property_id = P.kind_property_id
				LEFT JOIN cats_properties CP ON CP.property_id = P.property_id
			WHERE CP.cat_id = ?
			ORDER BY CP.pos ASC`

	if err := server.Db.Debug().Raw(query, catId).Scan(&list).Error; err != nil {
		return list, err
	}

	if err := valuesPropertyService.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}
func (ps PropertyService) GetPropertiesFullByCatIds(catIds []uint64) ([]*response.PropertyFull, error) {
	valuesPropertyService := NewValuesPropertyService()
	list := make([]*response.PropertyFull, 0)
	query := `
		SELECT  ` + ps.tblFields + `,
				KP.name AS kind_property_name, 
				CP.pos AS property_pos,
				CP.is_require AS property_is_require
			FROM properties P
				LEFT JOIN kind_properties KP ON KP.kind_property_id = P.kind_property_id
				LEFT JOIN cats_properties CP ON CP.property_id = P.property_id
			WHERE CP.cat_id IN (?)
			ORDER BY CP.pos ASC`

	if err := server.Db.Debug().Raw(query, catIds).Scan(&list).Error; err != nil {
		return list, err
	}

	if err := valuesPropertyService.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}
func (ps PropertyService) GetPropertyById(propertyId uint64) (*storage.Property, error) {
	property := new(storage.Property)

	return property, server.Db.Debug().First(property, propertyId).Error // проверяется в контроллере
}
func (ps PropertyService) GetPropertyFullById(propertyId uint64) (*response.PropertyFull, error) {
	propertyFull := new(response.PropertyFull)
	query := `
		SELECT ` + ps.tblFields + `, 
				KP.name AS kind_property_name 
			FROM properties P
				LEFT JOIN kind_properties KP ON KP.kind_property_id = P.kind_property_id
			WHERE P.property_id = ?`
	err := server.Db.Debug().Raw(query, propertyId).Scan(&propertyFull).Error // проверяется в контроллере

	if !gorm.IsRecordNotFoundError(err) {
		valuesPropertyService := NewValuesPropertyService()
		pValues, err := valuesPropertyService.GetValuesByPropertyId(propertyFull.PropertyId)
		if err != nil {
			return propertyFull, err
		}

		propertyFull.Values = pValues
	}

	return propertyFull, err
}
func (ps PropertyService) Create(property *storage.Property) error {
	if !server.Db.Debug().NewRecord(property) {
		return errOnNewRecordNewProperty
	}

	return server.Db.Debug().Create(property).Error
}
func (ps PropertyService) Update(property *storage.Property) error {
	return server.Db.Debug().Save(property).Error
}
func (ps PropertyService) Delete(propertyId uint64) error {
	property := storage.Property{
		PropertyId: propertyId,
	}

	return server.Db.Debug().Delete(property).Error
}
func (ps PropertyService) ReWriteValuesForProperties(
	propertyId uint64,
	mId map[string]string,
	mTitle map[string]string,
	mPos map[string]string) ([]storage.ValueProperty, error) {

	listOld := make([]storage.ValueProperty, 0)
	listNew := make([]storage.ValueProperty, 0)
	listResult := make([]storage.ValueProperty, 0)

	// возьмем старый список
	if err := server.Db.Debug().Where("property_id = ?", propertyId).Find(&listOld).Error; err != nil {
		return listResult, err
	}

	// создадим тут пришедший список
	for k, title := range mTitle {
		valueProperty := storage.ValueProperty{}
		valueProperty.Title = strings.TrimSpace(title)
		valueProperty.PropertyId = propertyId

		if val, found := mPos[k]; found {
			if iPos, err := strconv.ParseUint(val, 10, 64); err == nil && iPos > 0 {
				valueProperty.Pos = iPos
			}
		}
		if val, found := mId[k]; found {
			if iId, err := strconv.ParseUint(val, 10, 64); err == nil && iId > 0 {
				valueProperty.ValueId = iId
			}
		}

		listNew = append(listNew, valueProperty)
	}
	// \

	// возьмем id-шники которые надо удалить и при необходимости удалим
	removeId := make([]uint64, 0)
	for _, v1 := range listOld {
		isFound := false
		for _, v2 := range listNew {
			if v2.ValueId == v1.ValueId {
				isFound = true
			}
		}

		if !isFound {
			removeId = append(removeId, v1.ValueId)
		}
	}
	if len(removeId) > 0 {
		err := server.Db.Debug().Where("value_id IN (?)", removeId).Delete(storage.ValueProperty{}).Error
		if err != nil {
			return listResult, err
		}
	}
	// \

	// обновим/добавим остальные эл-ты
	for _, v := range listNew {
		if v.ValueId == 0 {
			if !server.Db.Debug().NewRecord(v) {
				return listResult, errOnNewRecordNewProperty
			}

			if err := server.Db.Debug().Create(&v).Error; err != nil {
				return listResult, err
			}

		} else {
			if err := server.Db.Debug().Model(&v).Update(v).Error; err != nil {
				return listResult, err
			}
		}

		listResult = append(listResult, v)
	}

	return listResult, nil
}

// private -------------------------------------------------------------------------------------------------------------
