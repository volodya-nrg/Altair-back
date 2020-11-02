package manager

import (
	"altair/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// InArray - находится ли значение в массиве
func InArray(val, array interface{}) (isFind bool, index int) {
	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				isFind = true
				index = i
				return
			}
		}
	}

	return
}

// RandStringRunes - случайная строка
func RandStringRunes(n int) string {
	return randGenerationString([]rune("abcdefghijklmnopqrstuvwxyz"), n)
}

// RandASCII - случайная строка и символов ASCII
func RandASCII(n int) string {
	return randGenerationString([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"), n)
}

// ValidateEmail - валидация е-мэйла
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return Re.MatchString(email)
}

// HashAndSalt - шифрование строки
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	if err != nil {
		logger.Error.Fatalln(err.Error())
	}

	return string(hash)
}

// ComparePasswords - сверка паролей (хеш и строка)
func ComparePasswords(hashedPwd, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		logger.Warning.Println(err.Error())
		return false
	}

	return true
}

// UploadImage - загрузка изображения
func UploadImage(file *multipart.FileHeader, pathUpload string, funcSave func(file *multipart.FileHeader, filePath string) error) (string, error) {
	var isPng bool
	var result string
	contentType := file.Header.Get("Content-Type")
	extension := ".jpg"

	// если заголовки правильно не выставлены на картинки, то Content-Type = application/octet-stream
	if exists, _ := InArray(contentType, []string{"image/jpeg", "image/png"}); !exists {
		return result, ErrNotAllowedImageType
	}
	if file.Size < 1 {
		return result, ErrFileSizeIsZero
	}
	if contentType == "image/png" {
		isPng = true
		extension = ".png"
	}

	fileHandler, err := file.Open()
	if err != nil {
		return result, fmt.Errorf("1, %s", err)
	}
	defer fileHandler.Close()

	result = fmt.Sprintf("%d%s%s", time.Now().Unix(), RandStringRunes(4), extension)
	myFilepath := pathUpload + "/" + result

	if err := funcSave(file, myFilepath); err != nil {
		return result, fmt.Errorf("3, %s", err)
	}

	imgFile, err := os.Open(myFilepath)
	if err != nil {
		return result, fmt.Errorf("4, %s", err)
	}
	defer imgFile.Close()

	imgCfg, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return result, fmt.Errorf("5, %s", err)
	}

	if imgCfg.Width < 1000 {
		return result, nil
	}

	// decode into image.Image
	var img image.Image

	buf, err := ioutil.ReadFile(imgFile.Name())
	if err != nil {
		return result, fmt.Errorf("6, %s", err)
	}

	pReaderImg := bytes.NewReader(buf)

	if isPng {
		img, err = png.Decode(pReaderImg)

	} else {
		img, err = jpeg.Decode(pReaderImg)
	}

	if err != nil {
		return result, fmt.Errorf("7, %s", err)
	}

	// resize to width 1000 using Lanczos resampling and preserve aspect ratio
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create(myFilepath)
	if err != nil {
		return result, fmt.Errorf("8, %s", err)
	}
	defer out.Close()

	// write new image to file
	if isPng {
		err = png.Encode(out, m)

	} else {
		err = jpeg.Encode(out, m, nil)
	}

	if err != nil {
		return result, fmt.Errorf("9, %s", err)
	}

	return result, nil
}

// FolderOrFileExists - проверка: папка или файл
func FolderOrFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// TranslitRuToEn - транслитерация
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

// GetYoutubeHash - достает хеш-код из ссылки Ютуба
func GetYoutubeHash(urlStr string) string {
	var result string

	if strings.HasPrefix(urlStr, "https://www.youtube.com/embed/") {
		return strings.ReplaceAll(urlStr, "https://www.youtube.com/embed/", "")
	}

	u, err := url.Parse(strings.TrimSpace(urlStr))
	if err != nil {
		return result
	}

	if u.RawQuery != "" {
		aQuery := strings.Split(u.RawQuery, "&")

		for _, v := range aQuery {
			aEq := strings.Split(v, "=")

			if len(aEq) > 1 {
				if aEq[0] == "v" {
					result = aEq[1]
				}
			}
		}

	} else if u.Path != "" {
		path := u.Path[1:]
		aPath := strings.Split(path, "/")

		if len(aPath) > 0 {
			result = aPath[0]
		}
	}

	return result
}

// RandIntByRange - сеет случайное число, от минимума до максимума
func RandIntByRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// IsSocialEmail - проверка на соц. е-мэйл
func IsSocialEmail(email string) bool {
	result := false
	aDomains := []string{"@vk.com", "@ok.ru", "@facebook.com", "@google.com"}

	for _, v := range aDomains {
		if strings.HasSuffix(email, v) {
			result = true
			break
		}
	}

	return result
}

// SToUint64 - конфертация строки в uint64
func SToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// MakeRequest - запрос GET, POST на удаленную машину
func MakeRequest(method, urlStr string, receiver interface{}, formData map[string]string) error {
	var err error
	var resp = new(http.Response)
	urlData := url.Values{}

	for k, v := range formData {
		urlData.Add(k, v)
	}

	if method == "get" {
		if len(urlData) > 0 {
			urlStr += "?" + urlData.Encode()
		}

		resp, err = http.Get(urlStr)

	} else {
		resp, err = http.PostForm(urlStr, urlData)
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if receiver != nil {
		if err := json.NewDecoder(resp.Body).Decode(receiver); err != nil {
			return err
		}
	}

	return nil
}

// PrettyPrint - красивый вывод данных
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")

	fmt.Println(string(s))
}

// GetUserIDAndRole - достает данные о текущем пользователе (userID и роль)
func GetUserIDAndRole(c *gin.Context) (userID uint64, userRole string, outError error) {
	var ok bool

	userID, ok = c.MustGet("userID").(uint64)
	if !ok {
		outError = ErrUndefinedUserID
		return
	}

	userRole, ok = c.MustGet("userRole").(string)
	if !ok {
		outError = ErrUndefinedUserRole
		return
	}

	return
}

// Lorem - сеет случайные слова
func Lorem(size int) string {
	const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`

	if size >= len(loremString) || size < 1 {
		return loremString
	}

	res := loremString[:size]
	res = strings.TrimSpace(res)

	return res
}

// GetTagsFromStruct - достает тег из структуры
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
func randGenerationString(letterRunes []rune, myLen int) string {
	rand.Seed(time.Now().UnixNano()) // сеет превдослучайное число (rand.Seed(time.Now().UTC().UnixNano()))

	b := make([]rune, myLen)
	l := len(letterRunes)

	for i := range b {
		b[i] = letterRunes[rand.Intn(l)]
	}

	return string(b)
}

// Вычисление времяни работы программы
// timeStart := time.Now()
// for i := 1000000; i >= 0; i-- {
//   Fn()
// }
// fmt.Println("====>", time.Since(timeStart))
