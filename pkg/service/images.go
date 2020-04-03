package service

import (
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"fmt"
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
	pImages := make([]*storage.Image, 0)
	err := server.Db.Debug().Where("el_id IN (?) AND opt = ?", elIds, opt).Find(&pImages).Error

	return pImages, err
}
func (is ImageService) GetImageById(imgId uint64) (*storage.Image, error) {
	img := new(storage.Image)
	err := server.Db.Debug().First(img, imgId).Error // проверяется в контроллере

	return img, err
}
func (is ImageService) Create(image *storage.Image) error {
	if !server.Db.Debug().NewRecord(image) {
		return errNotCreateNewImage
	}

	err := server.Db.Debug().Create(image).Error

	return err
}
func (is ImageService) Update(img *storage.Image) error {
	err := server.Db.Debug().Save(img).Error
	return err
}
func (is ImageService) Delete(imgId uint64) error {
	image, err := is.GetImageById(imgId)
	if err != nil {
		return err
	}

	if err := server.Db.Debug().Delete(storage.Image{}, "img_id = ?", image.ImgId).Error; err != nil {
		return err
	}

	myFilepath := fmt.Sprintf("%s%s", is.relativeImgDir, image.Filepath)
	if has := helpers.FileExists(myFilepath); has == true {
		_ = os.Remove(myFilepath)
	}

	return nil
}
func (is ImageService) DeleteAll(images []*storage.Image) error {
	for _, v := range images {
		if err := is.Delete(v.ImgId); err != nil {
			return err
		}
	}

	return nil
}
