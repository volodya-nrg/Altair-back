package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

// NewValuesPropService - фабрика, создает значения свойства
func NewValuesPropService() *ValuesPropService {
	return new(ValuesPropService)
}

// ValuesPropService - структура значений свойства
type ValuesPropService struct{}

// PopulateWithValues - обогатить значениями
func (vs ValuesPropService) PopulateWithValues(listProps []*response.PropFull) error {
	if len(listProps) < 1 {
		return nil
	}

	values := make([]*storage.ValueProp, 0)
	propIDs := make([]uint64, 0)

	for _, v := range listProps {
		propIDs = append(propIDs, v.PropID)
	}

	err := server.Db.Order("pos asc").Where("prop_id IN (?)", propIDs).Find(&values).Error
	if err != nil {
		return err
	}

	for _, prop := range listProps {
		for _, value := range values {
			if prop.PropID == value.PropID {
				prop.Values = append(prop.Values, value)
			}
		}
	}

	return nil
}

// GetValuesByPropID - получить значения относительно ID свойства
func (vs ValuesPropService) GetValuesByPropID(propID uint64) ([]*storage.ValueProp, error) {
	values := make([]*storage.ValueProp, 0)
	err := server.Db.Order("pos asc").Where("prop_id = ?", propID).Find(&values).Error

	return values, err
}

// DeleteByPropID - удалить относительно от ID свойства
func (vs ValuesPropService) DeleteByPropID(propID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Delete(storage.ValueProp{}, "prop_id = ?", propID).Error

	return err
}
