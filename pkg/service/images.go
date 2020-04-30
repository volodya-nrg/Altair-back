package service

import (
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

func NewImageService() *ImageService {
	imgS := new(ImageService)

	imgS.relativeImgDir = "./web/images/"

	return imgS
}

type ImageService struct {
	relativeImgDir string
}

func (is ImageService) GetImagesByElIdsAndOpt(elIds []uint64, opt string) ([]*storage.Image, error) {
	images := make([]*storage.Image, 0)
	err := server.Db.Debug().Where("el_id IN (?) AND opt = ?", elIds, opt).Find(&images).Error

	return images, err
}
func (is ImageService) GetImageById(imgId uint64) (*storage.Image, error) {
	img := new(storage.Image)
	err := server.Db.Debug().First(img, imgId).Error // проверяется в контроллере

	return img, err
}
func (is ImageService) Create(image *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	if !tx.NewRecord(image) {
		return errNotCreateNewImage
	}

	err := tx.Create(image).Error

	return err
}
func (is ImageService) Update(img *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(img).Error
	return err
}
func (is ImageService) Delete(img *storage.Image, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	if err := tx.Delete(img).Error; err != nil {
		return err
	}

	myFilepath := fmt.Sprintf("%s%s", is.relativeImgDir, img.Filepath)
	if has := helpers.FileExists(myFilepath); has == true {
		_ = os.Remove(myFilepath)
	}

	return nil
}
