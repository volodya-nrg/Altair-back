package service

import (
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

func NewKindPropService() *KindPropService {
	return new(KindPropService)
}

type KindPropService struct{}

func (ks KindPropService) GetKindProps(orderBy string) ([]*storage.KindProp, error) {
	props := make([]*storage.KindProp, 0)
	err := server.Db.Debug().Order(orderBy).Find(&props).Error // сортировка нужна

	return props, err
}
func (ks KindPropService) GetKindPropById(kindPropId uint64) (*storage.KindProp, error) {
	prop := new(storage.KindProp)
	err := server.Db.Debug().First(prop, kindPropId).Error

	return prop, err
}
func (ks KindPropService) Create(prop *storage.KindProp, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}
	if !server.Db.Debug().NewRecord(prop) {
		return errOnNewRecordNewKindProp
	}

	err := tx.Create(prop).Error

	return err
}
func (ks KindPropService) Update(prop *storage.KindProp, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(prop).Error

	return err
}
func (ks KindPropService) Delete(kindPropId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Delete(storage.KindProp{}, "kind_prop_id = ?", kindPropId).Error

	return err
}
