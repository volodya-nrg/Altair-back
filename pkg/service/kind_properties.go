package service

import (
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
)

func NewKindPropertyService() *KindPropertyService {
	return new(KindPropertyService)
}

type KindPropertyService struct{}

func (ks KindPropertyService) GetKindProperties() ([]*storage.KindProperty, error) {
	properties := make([]*storage.KindProperty, 0)
	err := server.Db.Debug().Order("kind_property_id asc").Find(&properties).Error // сортировка нужна для теста

	return properties, err
}
func (ks KindPropertyService) GetKindPropertyById(kindPropertyId uint64) (*storage.KindProperty, error) {
	property := new(storage.KindProperty)
	err := server.Db.Debug().First(property, kindPropertyId).Error

	return property, err
}
func (ks KindPropertyService) Create(property *storage.KindProperty, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}
	if !server.Db.Debug().NewRecord(property) {
		return errOnNewRecordNewKindProperty
	}

	err := tx.Create(property).Error

	return err
}
func (ks KindPropertyService) Update(property *storage.KindProperty, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(property).Error

	return err
}
func (ks KindPropertyService) Delete(kindPropertyId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Delete(storage.KindProperty{}, "kind_property_id = ?", kindPropertyId).Error

	return err
}
