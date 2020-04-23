package helpers

import (
	"altair/pkg/logger"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	errNotAllowedImageType = errors.New("not allowed image type")
	errFileSizeIsZero      = errors.New("file size is zero")
)

func InArray(val interface{}, array interface{}) (bool, int) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				return true, i
			}
		}
	}

	return false, -1
}
func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	// rand.Seed(time.Now().UTC().UnixNano()) // сеет превдослучайное число

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	l := len(letterRunes)

	for i := range b {
		b[i] = letterRunes[rand.Intn(l)]
	}

	return string(b)
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return Re.MatchString(email)
}
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	if err != nil {
		logger.Error.Fatalln(err)
	}

	return string(hash)
}
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		logger.Warning.Println(err)
		return false
	}

	return true
}
func UploadImage(file *multipart.FileHeader, pathUpload string, funcSave func(file *multipart.FileHeader, filePath string) error) (string, error) {
	contentType := file.Header.Get("Content-Type")
	isPng := false
	result := ""
	extension := ".jpg"

	if exists, _ := InArray(contentType, []string{"image/jpeg", "image/png"}); !exists {
		return "", errNotAllowedImageType
	}
	if file.Size < 1 {
		return "", errFileSizeIsZero
	}
	if contentType == "image/png" {
		isPng = true
		extension = ".png"
	}

	fileHandler, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("1, %s", err)
	}
	defer func() {
		_ = fileHandler.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, fileHandler); err != nil {
		return result, fmt.Errorf("2, %s", err)
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	filename := fmt.Sprintf("%s%d%s%s%s", hash[:4], time.Now().Unix(), "_", RandStringRunes(4), extension)

	for k, v := range filename {
		ch := rune(v)

		if k == 2 || k == 4 {
			result += "/"
		}

		result += string(ch)
	}
	myFilepath := pathUpload + "/" + result

	if err := os.MkdirAll(filepath.Dir(pathUpload+"/"+result), os.ModePerm); err != nil {
		return "", fmt.Errorf("10, %s", err)
	}

	if err := funcSave(file, myFilepath); err != nil {
		return "", fmt.Errorf("3, %s", err)
	}

	imgFile, err := os.Open(myFilepath)
	if err != nil {
		return "", fmt.Errorf("4, %s", err)
	}
	defer func() {
		_ = imgFile.Close()
	}()

	imgCfg, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return "", fmt.Errorf("5, %s", err)
	}

	if imgCfg.Width < 1000 {
		return result, nil
	}

	// decode into image.Image
	var img image.Image

	buf, err := ioutil.ReadFile(imgFile.Name())
	if err != nil {
		return "", fmt.Errorf("6, %s", err)
	}

	pReaderImg := bytes.NewReader(buf)

	if isPng {
		img, err = png.Decode(pReaderImg)

	} else {
		img, err = jpeg.Decode(pReaderImg)
	}

	if err != nil {
		return "", fmt.Errorf("7, %s", err)
	}

	// resize to width 1000 using Lanczos resampling and preserve aspect ratio
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create(myFilepath)
	if err != nil {
		return "", fmt.Errorf("8, %s", err)
	}
	defer func() {
		_ = out.Close()
	}()

	// write new image to file
	if isPng {
		err = png.Encode(out, m)

	} else {
		err = jpeg.Encode(out, m, nil)
	}

	if err != nil {
		return "", fmt.Errorf("9, %s", err)
	}

	return result, nil
}
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
func TranslitRuToEn(str string) string {
	chars := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "е": "e", "ё": "e", "ж": "zh", "з": "z", "и": "i", "й": "i",
		"к": "k", "л": "l", "м": "m", "н": "n", "о": "o", "п": "p", "р": "r", "с": "s", "т": "t", "у": "u", "ф": "f",
		"х": "kh", "ц": "ts", "ч": "ch", "ш": "sh", "щ": "shch", "ъ": "ie", "ы": "y", "э": "e", "ю": "iu", "я": "ia",
	}
	alphabetEn := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "k", "l", "m", "n", "o", "p", "q", "r", "s",
		"t", "v", "x", "y", "z"}
	result := ""

	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, "—", "-")
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.ReplaceAll(str, "_", "-")

	for _, v := range str {
		symbol := string(v)

		if unicode.IsDigit(v) || symbol == "-" {
			result += symbol
			continue
		}

		isFindEnChar := false
		for _, charEn := range alphabetEn {
			if charEn == symbol {
				result += charEn
				isFindEnChar = true
				break
			}
		}
		if isFindEnChar {
			continue
		}

		if v2, ok2 := chars[symbol]; ok2 {
			result += v2
			continue
		}
	}

	return result
}
func MakeRequest(method string, url string, result interface{}, formData ...map[string]interface{}) error {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	bytesRepresentation, err := json.Marshal(formData[0])
	if err != nil {
		logger.Warning.Println("-----> err in Marshal")
		return err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		logger.Warning.Println("-----> err in NewRequest")
		return err
	}

	if method == "post" {
		request.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(request)
	if err != nil {
		logger.Warning.Println("-----> err in client.Do")
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	logger.Warning.Println("-----> err in ioutil.ReadAll")
	//	return err
	//}
	//result = body

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Warning.Println("-----> Decode", result, err.Error())
		return err
	}
	PrettyPrint(result)

	return nil
}
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")

	fmt.Println(string(s))
}
func GetTagsFromStruct(i interface{}, tagName string) []string {
	aTagName := make([]string, 0)

	v := reflect.ValueOf(i)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		aTagName = append(aTagName, f.Tag.Get(tagName))
	}

	return aTagName
}

// private -------------------------------------------------------------------------------------------------------------
// Вычисление времяни работы программы
// timeStart := time.Now()
// for i := 1000000; i >= 0; i-- {
//   Fn()
// }
// fmt.Println("====>", time.Since(timeStart))
