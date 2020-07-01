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

	//check, err := ms.UploadCheck(fileInfo.Name(), hash, fileInfo.Size())
	//if err != nil {
	//	return err
	//}
	//if check.FileExists == "yes" || check.HashExists == "yes" {
	//	return errors.New("file/hash already exists")
	//}

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

	// spew.Dump(req)

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

//func (ms *MediafireService) UserRenewSessionToken() error {
//	urlPath := "/api/1.5/user/renew_session_token.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"response_format": "json",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUserRenewSessionToken
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	ms.Data.Response.SessionToken = receiver.Response.SessionToken
//
//	return nil
//}
//func (ms *MediafireService) UserGetActionToken() (string, error) {
//	var actionToken string
//	urlPath := "/api/1.5/user/get_action_token.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"type":            "upload",
//		"response_format": "json",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return actionToken, err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUserGetActionToken
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return actionToken, err
//	}
//
//	return receiver.Response.ActionToken, nil
//}
//func (ms *MediafireService) UserDestroyActionToken(actionToken string) error {
//	urlPath := "/api/1.5/user/destroy_action_token.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"action_token":    actionToken,
//		"response_format": "json",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUserDestroyActionToken
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//func (ms *MediafireService) UserGetInfo() error {
//	urlPath := "/api/1.5/user/get_info.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"response_format": "json",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUserGetInfo
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (ms *MediafireService) UploadGetOptions() error {
//	urlPath := "/api/1.5/upload/get_options.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"response_format": "json",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUploadGetOptions
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//func (ms *MediafireService) UploadSetOptions() error {
//	urlPath := "/api/1.5/upload/set_options.php"
//	query := map[string]string{
//		"session_token":       ms.Data.Response.SessionToken,
//		"response_format":     "json",
//		"action_on_duplicate": "keep",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUploadSetOptions
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//func (ms *MediafireService) UploadCheck(filename string, hash string, size int64) (MediafireUploadCheckResponse, error) {
//	var result MediafireUploadCheckResponse
//	urlPath := "/api/1.5/upload/check.php"
//	query := map[string]string{
//		"filename":        filename,
//		"hash":            hash,
//		"size":            fmt.Sprint(size),
//		"session_token":   ms.Data.Response.SessionToken,
//		"response_format": "json",
//		"folder_key":      ms.FolderKey,
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return result, err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireUploadCheck
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return result, err
//	}
//
//	availableSpace, err := manager.SToUint64(receiver.Response.AvailableSpace)
//	if err != nil {
//		return receiver.Response, err
//	}
//
//	if availableSpace < size {
//		return receiver.Response, fmt.Errorf("available space is low (%d)", availableSpace)
//	}
//
//	return receiver.Response, nil
//}
//
//func (ms *MediafireService) FolderConfigureFiledrop() (MediafireFolderConfigureFiledropResponse, error) {
//	var result MediafireFolderConfigureFiledropResponse
//	urlPath := "/api/1.5/folder/configure_filedrop.php"
//	query := map[string]string{
//		"session_token":   ms.Data.Response.SessionToken,
//		"response_format": "json",
//		"action":          "disable",
//		"folder_key":      ms.FolderKey,
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return result, err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireFolderConfigureFiledrop
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return result, err
//	}
//
//	return receiver.Response, nil
//}
//func (ms *MediafireService) FolderGetInfo() error {
//	urlPath := "/api/1.5/folder/get_info.php"
//	query := map[string]string{
//		"response_format": "json",
//		"folder_key":      ms.FolderKey,
//		"session_token":   ms.Data.Response.SessionToken,
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireFolderGetInfo
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//func (ms *MediafireService) FolderGetContent() error {
//	// чет не работает
//	urlPath := "/api/1.5/folder/get_content.php"
//	query := map[string]string{
//		"response_format": "json",
//		"folder_key":      ms.FolderKey,
//		//"folder_path": "altair",
//		"session_token": ms.Data.Response.SessionToken,
//		// "content_type":    "files",
//	}
//	urlData := ms.collectURLData(query, urlPath, true)
//
//	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
//	// resp, err := http.Get(ms.Domain+urlPath + "?" + urlData.Encode())
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	var receiver MediafireFolderGetInfo
//	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
//		return err
//	}
//
//	return nil
//}
//func (ms *MediafireService) HasFolder(name string) (bool, error) {
//	var has bool
//	return has, nil
//}

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
func (ms *MFService) createSignature(url string) string {
	str := fmt.Sprintf("%d%s%s", ms.Data.Response.SecretKeyInt%256, ms.Data.Response.Time, url)
	signature := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return signature
}
