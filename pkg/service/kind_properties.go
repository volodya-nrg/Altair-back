package service

import (
	"altair/server"
	"altair/storage"
	"errors"
)

var (
	errOnNewRecordNewKindProperty = errors.New("err on NewRecord new kindProperty")
)

func NewKindPropertyService() *KindPropertyService {
	return new(KindPropertyService)
}

type KindPropertyService struct{}

func (ks KindPropertyService) GetKindProperties() ([]*storage.KindProperty, error) {
	properties := make([]*storage.KindProperty, 0)
	err := server.Db.Debug().Order("kind_property_id", true).Find(properties).Error // сортировка нужна для теста

	return properties, err
}
func (ks KindPropertyService) GetKindPropertyById(kindPropertyId uint64) (*storage.KindProperty, error) {
	property := new(storage.KindProperty)
	err := server.Db.Debug().First(property, kindPropertyId).Error // на notFound проверяется в контроллере

	return property, err
}
func (ks KindPropertyService) Create(property *storage.KindProperty) error {
	if !server.Db.Debug().NewRecord(property) {
		return errOnNewRecordNewKindProperty
	}

	return server.Db.Debug().Create(property).Error
}
func (ks KindPropertyService) Update(property *storage.KindProperty) error {
	return server.Db.Debug().Save(property).Error
}
func (ks KindPropertyService) Delete(kindPropertyId uint64) error {
	property := storage.KindProperty{
		KindPropertyId: kindPropertyId,
	}

	return server.Db.Debug().Delete(property).Error
}

// private -------------------------------------------------------------------------------------------------------------
