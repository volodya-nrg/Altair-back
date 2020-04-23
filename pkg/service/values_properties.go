package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

func NewValuesPropertyService() *ValuesPropertyService {
	return new(ValuesPropertyService)
}

type ValuesPropertyService struct{}

func (vs ValuesPropertyService) PopulateWithValues(listProperties []*response.PropertyFull) error {
	if len(listProperties) < 1 {
		return nil
	}

	values := make([]*storage.ValueProperty, 0)
	propIds := make([]uint64, 0)

	for _, v := range listProperties {
		propIds = append(propIds, v.PropertyId)
	}

	err := server.Db.Debug().
		Order("pos asc").
		Where("property_id IN (?)", propIds).
		Find(&values).Error
	if err != nil {
		return err
	}

	for _, property := range listProperties {
		for _, value := range values {
			if property.PropertyId == value.PropertyId {
				property.Values = append(property.Values, value)
			}
		}
	}

	return nil
}
func (vs ValuesPropertyService) GetValuesByPropertyId(propId uint64) ([]*storage.ValueProperty, error) {
	values := make([]*storage.ValueProperty, 0)
	err := server.Db.Debug().
		Order("pos asc").
		Where("property_id = ?", propId).
		Find(&values).Error

	return values, err
}
func (vs ValuesPropertyService) DeleteByPropertyId(propertyId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Delete(storage.ValueProperty{}, "property_id = ?", propertyId).Error

	return err
}
