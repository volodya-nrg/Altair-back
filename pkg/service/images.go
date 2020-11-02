package service

import (
	"altair/server"
	"altair/storage"
	"gorm.io/gorm"
)

// NewImageService - фабрика, создает объект Изображения
func NewImageService() *ImageService {
	return new(ImageService)
}

// ImageService - структура изображения
type ImageService struct{}

// GetImagesByElIDsAndOpt - получить картинки относительно ID элемента и опции
func (is ImageService) GetImagesByElIDsAndOpt(elIDs []uint64, opt string) ([]*storage.Image, error) {
	images := make([]*storage.Image, 0)
	err := server.Db.Where("opt = ? AND el_id IN (?)", opt, elIDs).Find(&images).Error

	return images, err
}

// GetImageByID - получить изображения относительно его ID
func (is ImageService) GetImageByID(imgID uint64) (*storage.Image, error) {
	img := new(storage.Image)
	err := server.Db.First(img, imgID).Error // проверяется в контроллере

	return img, err
}

// Create - создать запись об изображении
func (is ImageService) Create(image *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(image).Error

	return err
}

// Update - изменить запись об изображении
func (is ImageService) Update(img *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(img).Error
	return err
}

// Delete - удалить запись об изображении
func (is ImageService) Delete(img *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	if err := tx.Delete(img).Error; err != nil {
		return err
	}

	// тут по идеи надо удалить файлы с удаленного сервера

	return nil
}
