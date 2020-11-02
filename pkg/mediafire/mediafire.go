package mediafire

import (
	"altair/configs"
	"altair/pkg/manager"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// NewMediafireService - фабрика, создает объект Медиафайр
func NewMediafireService() *MFService {
	mf := new(MFService)

	mf.AppID = configs.Cfg.Mediafire.AppID
	mf.AppName = configs.Cfg.Mediafire.AppName
	mf.APIKey = configs.Cfg.Mediafire.APIKey
	mf.UserEmail = configs.Cfg.Mediafire.UserEmail
	mf.UserPassword = configs.Cfg.Mediafire.UserPassword
	mf.Domain = configs.Cfg.Mediafire.Domain
	mf.FolderKey = configs.Cfg.Mediafire.FolderKey

	return mf
}

// MFService - структура Медиафайр
type MFService struct {
	AppID        string
	AppName      string
	APIKey       string
	UserEmail    string
	UserPassword string
	Domain       string
	FolderKey    string
	Data         SessionToken
}

// UserGetSessionToken - получить сессионный токен
func (ms *MFService) UserGetSessionToken() error {
	urlPath := "/api/1.5/user/get_session_token.php"
	signature := sha1.Sum([]byte(fmt.Sprintf("%s%s%s%s", ms.UserEmail, ms.UserPassword, ms.AppID, ms.APIKey)))
	query := map[string]string{
		"response_format": "json",
		"email":           ms.UserEmail,
		"password":        ms.UserPassword,
		"application_id":  ms.AppID,
		"signature":       fmt.Sprintf("%x", signature),
		"token_version":   "2",
	}
	urlData := ms.collectURLData(query, urlPath, false)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := ms.afterResponse(&ms.Data, resp.Body); err != nil {
		return err
	}

	ms.Data.Response.SecretKeyInt, err = manager.SToUint64(ms.Data.Response.SecretKey)
	if err != nil {
		return err
	}

	return nil
}

// FolderCreate - создать папку
func (ms *MFService) FolderCreate(name string) (*FolderCreate, error) {
	urlPath := "/api/1.5/folder/create.php"
	query := map[string]string{
		"session_token":       ms.Data.Response.SessionToken,
		"foldername":          name,
		"parent_key":          ms.FolderKey,
		"action_on_duplicate": "replace",
		"response_format":     "json",
	}
	urlData := ms.collectURLData(query, urlPath, true)
	result := new(FolderCreate)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if err := ms.afterResponse(result, resp.Body); err != nil {
		return result, err
	}

	return result, nil
}

// UploadPollUpload - загрузка в режиме Poll
func (ms *MFService) UploadPollUpload(key string) (string, error) {
	var result string
	urlPath := "/api/1.5/upload/poll_upload.php"
	query := map[string]string{
		"key":             key,
		"response_format": "json",
	}
	urlData := ms.collectURLData(query, urlPath, false)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	receiver := new(UploadPollUpload)
	if err := json.NewDecoder(resp.Body).Decode(receiver); err != nil {
		return result, err
	}
	if receiver.Response.Result == "Success" {
		result = receiver.Response.Doupload.Quickkey + "/" + receiver.Response.Doupload.Filename
		// z89ypihvuecivzr/test4(14).jpg
	}

	return result, nil
}

// UploadSimple - простая загрузка файлов на удаленный сервер
func (ms *MFService) UploadSimple(filepath string) (string, error) {
	var result string
	urlPath := "/api/1.5/upload/simple.php"
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	if !manager.FolderOrFileExists(filepath) {
		return result, errors.New("not found file")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return result, err
	}

	var fileBytes []byte
	h := sha256.New()
	if _, err := file.Read(fileBytes); err != nil {
		return result, err
	}
	hash := fmt.Sprintf("%x", h.Sum(fileBytes))

	// CREATE FOLDER
	now := time.Now()
	day := now.Day()
	folderDir := now.AddDate(0, 0, -1*(day-1)).Format("2006-01-02")

	if err := ms.UserGetSessionToken(); err != nil {
		return result, err
	}

	folderKey, err := ms.FolderCreate(folderDir)
	if err != nil {
		return result, err
	}

	query := map[string]string{
		"response_format":     "json",
		"session_token":       ms.Data.Response.SessionToken,
		"action_on_duplicate": "replace",
		"folder_key":          folderKey.Response.FolderKey,
	}

	urlData := ms.collectURLData(query, urlPath, true)
	body := new(bytes.Buffer)
	multiPartWriter := multipart.NewWriter(body)

	part, err := multiPartWriter.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		return result, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return result, err
	}

	err = multiPartWriter.Close()
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(http.MethodPost, ms.Domain+urlPath+"?"+urlData.Encode(), body)
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	req.Header.Set("X-Filename", fileInfo.Name())
	req.Header.Set("X-Filesize", fmt.Sprint(fileInfo.Size()))
	req.Header.Set("X-Filehash", hash)

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	receiver := new(UploadSimple)
	if err := ms.afterResponse(receiver, resp.Body); err != nil {
		return result, err
	}
	if receiver.Response.Doupload.Key != "" {
		filepath, err := ms.UploadPollUpload(receiver.Response.Doupload.Key)
		if err != nil {
			return result, err
		}

		result = filepath
	}

	return result, nil
}

// private -------------------------------------------------------------------------------------------------------------
func (ms *MFService) collectURLData(query map[string]string, urlPath string, addSignature bool) url.Values {
	urlData := url.Values{}

	for k, v := range query {
		urlData.Add(k, v)
	}

	if urlPath != "" && addSignature {
		urlData.Add("signature", ms.createSignature(urlPath+"?"+urlData.Encode()))
	}

	return urlData
}
func (ms *MFService) afterResponse(receiver CommonInterface, body io.Reader) error {
	if err := json.NewDecoder(body).Decode(receiver); err != nil {
		return err
	}
	if receiver.getCommon().Result == "Error" {
		return fmt.Errorf("%s (%d)", receiver.getCommon().Message, receiver.getCommon().Error)
	}
	if receiver.getCommon().NewKey == "yes" {
		ms.generateNewKey()
	}

	return nil
}
func (ms *MFService) generateNewKey() {
	ms.Data.Response.SecretKeyInt = (ms.Data.Response.SecretKeyInt * 16807) % 2147483647
	ms.Data.Response.SecretKey = fmt.Sprint(ms.Data.Response.SecretKeyInt)
}
func (ms *MFService) createSignature(urlSrc string) string {
	str := fmt.Sprintf("%d%s%s", ms.Data.Response.SecretKeyInt%256, ms.Data.Response.Time, urlSrc)
	signature := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return signature
}
