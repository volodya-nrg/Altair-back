package service

import (
	"altair/server"
	"altair/storage"
	"gorm.io/gorm"
)

// NewKindPropService - фабрика, создает объект вида свойства
func NewKindPropService() *KindPropService {
	return new(KindPropService)
}

// KindPropService - структура вида свойства
type KindPropService struct{}

// GetKindProps - получить виды свойств
func (ks KindPropService) GetKindProps(orderBy string) ([]*storage.KindProp, error) {
	props := make([]*storage.KindProp, 0)
	err := server.Db.Order(orderBy).Find(&props).Error // сортировка нужна

	return props, err
}

// GetKindPropByID - получить вид свойства относительно его ID
func (ks KindPropService) GetKindPropByID(kindPropID uint64) (*storage.KindProp, error) {
	prop := new(storage.KindProp)
	err := server.Db.First(prop, kindPropID).Error

	return prop, err
}

// Create - создать вид свойства
func (ks KindPropService) Create(prop *storage.KindProp, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(prop).Error

	return err
}

// Update - изменить вид свойства
func (ks KindPropService) Update(prop *storage.KindProp, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(prop).Error

	return err
}

// Delete - удалить вид свойства
func (ks KindPropService) Delete(kindPropID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Delete(storage.KindProp{}, "kind_prop_id = ?", kindPropID).Error

	return err
}
