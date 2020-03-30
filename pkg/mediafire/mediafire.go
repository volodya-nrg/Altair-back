package mediafire

import (
	"altair/configs"
	"altair/pkg/helpers"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func NewMediafireService() *MediafireService {
	mf := new(MediafireService)

	mf.AppID = configs.Cfg.Mediafire.AppID
	mf.AppName = configs.Cfg.Mediafire.AppName
	mf.APIKey = configs.Cfg.Mediafire.APIKey
	mf.UserEmail = configs.Cfg.Mediafire.UserEmail
	mf.UserPassword = configs.Cfg.Mediafire.UserPassword
	mf.Domain = configs.Cfg.Mediafire.Domain
	mf.FolderKey = configs.Cfg.Mediafire.FolderKey

	return mf
}

type MediafireService struct {
	AppID        string
	AppName      string
	APIKey       string
	UserEmail    string
	UserPassword string
	Domain       string
	FolderKey    string
	Data         MediafireSessionToken
}

func (ms *MediafireService) UserGetSessionToken() error {
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
	urlData := ms.collectUrlData(query, urlPath, false)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := ms.afterResponse(&ms.Data, resp.Body); err != nil {
		return err
	}

	ms.Data.Response.SecretKeyInt, err = strconv.ParseUint(ms.Data.Response.SecretKey, 10, 64)
	if err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UserRenewSessionToken() error {
	urlPath := "/api/1.5/user/renew_session_token.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUserRenewSessionToken
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	ms.Data.Response.SessionToken = receiver.Response.SessionToken

	return nil
}
func (ms *MediafireService) UserGetActionToken() (string, error) {
	var actionToken string
	urlPath := "/api/1.5/user/get_action_token.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"type":            "upload",
		"response_format": "json",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return actionToken, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUserGetActionToken
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return actionToken, err
	}

	return receiver.Response.ActionToken, nil
}
func (ms *MediafireService) UserDestroyActionToken(actionToken string) error {
	urlPath := "/api/1.5/user/destroy_action_token.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"action_token":    actionToken,
		"response_format": "json",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUserDestroyActionToken
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UserGetInfo() error {
	urlPath := "/api/1.5/user/get_info.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUserGetInfo
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}

func (ms *MediafireService) UploadGetOptions() error {
	urlPath := "/api/1.5/upload/get_options.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUploadGetOptions
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UploadSetOptions() error {
	urlPath := "/api/1.5/upload/set_options.php"
	query := map[string]string{
		"session_token":       ms.Data.Response.SessionToken,
		"response_format":     "json",
		"action_on_duplicate": "keep",
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUploadSetOptions
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UploadAddWebUpload() error {
	urlPath := "/api/1.5/upload/add_web_upload.php"
	query := map[string]string{
		"response_format": "json",
		"session_token":   ms.Data.Response.SessionToken,
		"url":             "https://deswal.ru/wide/1920-1200/00000645.jpg",
		"filename":        "00000645.jpg",
		"folder_key":      ms.FolderKey,
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUploadAddWebUpload
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UploadCheck(filename string, hash string, size int64) (MediafireUploadCheckResponse, error) {
	var result MediafireUploadCheckResponse
	urlPath := "/api/1.5/upload/check.php"
	query := map[string]string{
		"filename":        filename,
		"hash":            hash,
		"size":            fmt.Sprint(size),
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
		"folder_key":      ms.FolderKey,
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return result, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUploadCheck
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return result, err
	}

	availableSpace, err := strconv.ParseInt(receiver.Response.AvailableSpace, 10, 64)
	if err != nil {
		return receiver.Response, err
	}

	if availableSpace < size {
		return receiver.Response, fmt.Errorf("available space is low (%d)", availableSpace)
	}

	return receiver.Response, nil
}
func (ms *MediafireService) UploadSimple(filepath string) error {
	urlPath := "/api/1.5/upload/simple.php"
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	if !helpers.FileExists(filepath) {
		return errors.New("not found file")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	sFileSize := fmt.Sprint(fileInfo.Size())

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	check, err := ms.UploadCheck(fileInfo.Name(), hash, fileInfo.Size())
	if err != nil {
		return err
	}

	if check.FileExists == "yes" || check.HashExists == "yes" {
		return errors.New("file/hash already exists")
	}

	actionToken, err := ms.UserGetActionToken()
	if err != nil {
		return err
	}
	defer func() {
		_ = ms.UserDestroyActionToken(actionToken)
	}()

	query := map[string]string{
		"response_format": "xml",       // пока только так
		"session_token":   actionToken, // ms.Data.Response.SessionToken,
		"folder_key":      ms.FolderKey,
	}
	urlData := ms.collectUrlData(query, urlPath, false)
	body := new(bytes.Buffer)
	multiPartWriter := multipart.NewWriter(body)

	for k, v := range query {
		if err := multiPartWriter.WriteField(k, v); err != nil {
			return err
		}
	}

	if err := multiPartWriter.WriteField("signature", ms.createSignature(urlPath+"?"+urlData.Encode())); err != nil {
		return err
	}

	part, err := multiPartWriter.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = multiPartWriter.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ms.Domain+urlPath, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType()) // multipart/form-data; boundary=7fb6da2fe7a1d4da1520382b5a878d76ff2af9838cdc4b4dc1ed5ecd0069
	req.Header.Set("x-filename", fileInfo.Name())
	req.Header.Set("x-filesize", sFileSize)
	req.Header.Set("x-filehash", hash)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(dataBytes))

	return nil
}
func (ms *MediafireService) UploadInstant(filepath string) error {
	urlPath := "/api/1.5/upload/instant.php"
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	if !helpers.FileExists(filepath) {
		return errors.New("not found file")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	//_, err = ms.UploadCheck(fileInfo.Name(), hash, fileInfo.Size())
	//if err != nil {
	//	return err
	//}

	//if err := ms.UserRenewSessionToken(); err != nil {
	//	return err
	//}

	//actionToken, err := ms.UserGetActionToken()
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	_ = ms.UserDestroyActionToken(actionToken)
	//}()

	//filedrop, err := ms.FolderConfigureFiledrop()
	//if err != nil {
	//	return err
	//}

	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
		"size":            fmt.Sprint(fileInfo.Size()),
		"hash":            hash,
		"filename":        fileInfo.Name(),
	}
	urlData := ms.collectUrlData(query, urlPath, false)
	body := new(bytes.Buffer)
	multiPartWriter := multipart.NewWriter(body)

	for k, v := range query {
		if err := multiPartWriter.WriteField(k, v); err != nil {
			return err
		}
	}

	if err := multiPartWriter.WriteField("signature", ms.createSignature(urlPath+"?"+urlData.Encode())); err != nil {
		return err
	}

	part, err := multiPartWriter.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = multiPartWriter.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ms.Domain+urlPath, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireUploadInstant
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}
func (ms *MediafireService) UploadResumable(filepath string) error {
	urlPath := "/api/1.5/upload/resumable.php"
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	if !helpers.FileExists(filepath) {
		return errors.New("not found file")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	sFileSize := fmt.Sprint(fileInfo.Size())

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	_, err = ms.UploadCheck(fileInfo.Name(), hash, fileInfo.Size())
	if err != nil {
		return err
	}

	//if err := ms.UserRenewSessionToken(); err != nil {
	//	return err
	//}

	//actionToken, err := ms.UserGetActionToken()
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	_ = ms.UserDestroyActionToken(actionToken)
	//}()

	query := map[string]string{
		"response_format": "xml",
		"session_token":   ms.Data.Response.SessionToken,
		"folder_key":      ms.FolderKey,
	}
	urlData := ms.collectUrlData(query, urlPath, false)
	body := new(bytes.Buffer)
	multiPartWriter := multipart.NewWriter(body)

	for k, v := range query {
		if err := multiPartWriter.WriteField(k, v); err != nil {
			return err
		}
	}

	if err := multiPartWriter.WriteField("signature", ms.createSignature(urlPath+"?"+urlData.Encode())); err != nil {
		return err
	}

	part, err := multiPartWriter.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = multiPartWriter.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ms.Domain+urlPath, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	req.Header.Set("x-filesize", sFileSize)
	req.Header.Set("x-filehash", hash)
	req.Header.Set("x-unit-hash", hash)
	req.Header.Set("x-unit-id", "0")
	req.Header.Set("x-unit-size", sFileSize)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(dataBytes))

	return nil
}

func (ms *MediafireService) FolderConfigureFiledrop() (MediafireFolderConfigureFiledropResponse, error) {
	var result MediafireFolderConfigureFiledropResponse
	urlPath := "/api/1.5/folder/configure_filedrop.php"
	query := map[string]string{
		"session_token":   ms.Data.Response.SessionToken,
		"response_format": "json",
		"action":          "disable",
		"folder_key":      ms.FolderKey,
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return result, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireFolderConfigureFiledrop
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return result, err
	}

	return receiver.Response, nil
}
func (ms *MediafireService) FolderGetInfo() error {
	urlPath := "/api/1.5/folder/get_info.php"
	query := map[string]string{
		"response_format": "json",
		"folder_key":      ms.FolderKey,
		"session_token":   ms.Data.Response.SessionToken,
	}
	urlData := ms.collectUrlData(query, urlPath, true)

	resp, err := http.PostForm(ms.Domain+urlPath, urlData)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var receiver MediafireFolderGetInfo
	if err := ms.afterResponse(&receiver, resp.Body); err != nil {
		return err
	}

	return nil
}

// private -------------------------------------------------------------------------------------------------------------
func (ms *MediafireService) afterResponse(receiver MediafireCommonInterface, body io.Reader) error {
	if err := json.NewDecoder(body).Decode(receiver); err != nil {
		return err
	}

	helpers.PrettyPrint(receiver)

	if receiver.getMediafireCommon().Result == "Error" {
		return fmt.Errorf("%s (%d)", receiver.getMediafireCommon().Message, receiver.getMediafireCommon().Error)
	}
	if receiver.getMediafireCommon().NewKey == "yes" {
		ms.generateNewKey()
	}

	return nil
}
func (ms *MediafireService) generateNewKey() {
	ms.Data.Response.SecretKeyInt = (ms.Data.Response.SecretKeyInt * 16807) % 2147483647
	ms.Data.Response.SecretKey = fmt.Sprint(ms.Data.Response.SecretKeyInt)
}
func (ms *MediafireService) collectUrlData(query map[string]string, urlPath string, addSignature bool) url.Values {
	urlData := url.Values{}

	for k, v := range query {
		urlData.Add(k, v)
	}

	if urlPath != "" && addSignature {
		urlData.Add("signature", ms.createSignature(urlPath+"?"+urlData.Encode()))
	}

	return urlData
}
func (ms *MediafireService) createSignature(url string) string {
	helpers.PrettyPrint(ms.Data.Response.SecretKeyInt)
	str := fmt.Sprintf("%d%s%s", ms.Data.Response.SecretKeyInt%256, ms.Data.Response.Time, url)

	helpers.PrettyPrint(str)

	signature := fmt.Sprintf("%x", md5.Sum([]byte(str)))

	return signature
}
