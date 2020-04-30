package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

func NewValuesPropService() *ValuesPropService {
	return new(ValuesPropService)
}

type ValuesPropService struct{}

func (vs ValuesPropService) PopulateWithValues(listProps []*response.PropFull) error {
	if len(listProps) < 1 {
		return nil
	}

	values := make([]*storage.ValueProp, 0)
	propIds := make([]uint64, 0)

	for _, v := range listProps {
		propIds = append(propIds, v.PropId)
	}

	err := server.Db.Debug().
		Order("pos asc").
		Where("prop_id IN (?)", propIds).
		Find(&values).Error
	if err != nil {
		return err
	}

	for _, prop := range listProps {
		for _, value := range values {
			if prop.PropId == value.PropId {
				prop.Values = append(prop.Values, value)
			}
		}
	}

	return nil
}
func (vs ValuesPropService) GetValuesByPropId(propId uint64) ([]*storage.ValueProp, error) {
	values := make([]*storage.ValueProp, 0)
	err := server.Db.Debug().
		Order("pos asc").
		Where("prop_id = ?", propId).
		Find(&values).Error

	return values, err
}
func (vs ValuesPropService) DeleteByPropId(propId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Delete(storage.ValueProp{}, "prop_id = ?", propId).Error

	return err
}
