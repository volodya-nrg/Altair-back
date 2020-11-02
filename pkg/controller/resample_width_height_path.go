package controller

import (
	"altair/pkg/logger"
	"altair/pkg/manager"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetResampleWidthHeightPath - изменение пропорций изображения (налету)
func GetResampleWidthHeightPath(c *gin.Context) {
	wSrc := c.Param("width")
	hSrc := c.Param("height")
	pathSrc := c.Param("path") // идет со слешем в начале

	if matched, err := regexp.Match(manager.ImgPathPattern, []byte(pathSrc)); err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return

	} else if !matched {
		c.JSON(404, manager.ErrNotMatched.Error())
		return
	}

	w, err := manager.SToUint64(wSrc)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	h, err := manager.SToUint64(hSrc)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(400, err.Error())
		return
	}

	path := pathSrc[1:]
	externalLink := "http://www.mediafire.com/file/" + path + "/file"
	fileNameOriginal := filepath.Base(path)
	fileNameNew := fmt.Sprintf("%dx%d_%s", w, h, fileNameOriginal)
	fileInDirResample := manager.DirResample + "/" + fileNameNew

	// отдаем существующий уже файл (с такой же шириной и высотой)
	if manager.FolderOrFileExists(fileInDirResample) {
		c.File(fileInDirResample)
		return
	}

	// считываем данные с внешного ресурса. Тут особоый случай, вспомог-ую ф-ию лучше не использовать
	resp, err := http.Get(externalLink)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Warning.Println(resp.Status)
		c.JSON(resp.StatusCode, resp.Status)
		return
	}

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	// если запрашиваемые параметры больше положенного, то отдать оригинал
	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(dataBytes))
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	if int(w) > imgCfg.Width || int(h) > imgCfg.Height {
		c.File(externalLink)
		return
	}

	isPng := false

	if strings.HasSuffix(fileNameOriginal, ".png") {
		isPng = true
	}

	var img image.Image
	pReaderImg := bytes.NewReader(dataBytes)

	if isPng {
		img, err = png.Decode(pReaderImg)

	} else {
		img, err = jpeg.Decode(pReaderImg)
	}

	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	var m image.Image

	switch {
	case w > 0 && h < 1:
		m = resize.Resize(uint(w), 0, img, resize.Lanczos3)
	case w < 1 && h > 0:
		m = resize.Resize(0, uint(h), img, resize.Lanczos3)
	default:
		m = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	}

	out, err := os.Create(fileInDirResample)
	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}
	defer out.Close()

	if isPng {
		err = png.Encode(out, m)

	} else {
		err = jpeg.Encode(out, m, nil)
	}

	if err != nil {
		logger.Warning.Println(err.Error())
		c.JSON(500, err.Error())
		return
	}

	c.File(fileInDirResample)
}
