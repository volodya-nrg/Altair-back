package mediafire

import "time"

// CommonInterface - общий интерфейс
type CommonInterface interface {
	getCommon() Common
}

// Common - общая структура данных
type Common struct {
	Error             int    `json:"error"`
	Message           string `json:"message"`
	Deprecated        string `json:"deprecated"`
	Result            string `json:"result"`
	NewKey            string `json:"new_key"`
	CurrentAPIVersion string `json:"current_api_version"`
	Action            string `json:"action"`
}

// SessionToken - структура сессии токена
type SessionToken struct {
	Response SessionTokenResponse `json:"response"`
}

// SessionTokenResponse - структура сессии токена (ответ)
type SessionTokenResponse struct {
	Common
	SessionToken string `json:"session_token"`
	SecretKey    string `json:"secret_key"`
	SecretKeyInt uint64
	Time         string `json:"time"`
	Ekey         string `json:"ekey"`
	Pkey         string `json:"pkey"`
}

func (m *SessionToken) getCommon() Common {
	return m.Response.Common
}

// UploadSimple - структура простой загрузки
type UploadSimple struct {
	Response UploadSimpleResponse `xml:"response" json:"response"`
}

// UploadSimpleResponse - структура ответа простой загрузки
type UploadSimpleResponse struct {
	Common
	Key               string                       `json:"key"`
	Server            string                       `json:"server"`
	NewDeviceRevision int                          `json:"new_device_revision"`
	Doupload          UploadSimpleResponseDoupload `json:"doupload"`
}

// UploadSimpleResponseDoupload - структура Doupload
type UploadSimpleResponseDoupload struct {
	Result string `json:"result"`
	Key    string `json:"key"`
}

func (m *UploadSimple) getCommon() Common {
	return m.Response.Common
}

//----------------------------------------------------------------------------------------------------------------------

// FolderCreate - структура создания папки
type FolderCreate struct {
	Response FolderCreateResponse `xml:"response" json:"response"`
}

// FolderCreateResponse - структура ответа создания папки
type FolderCreateResponse struct {
	Common
	FolderKey       string      `json:"folder_key"`
	ParentFolderkey string      `json:"parent_folderkey"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Created         string      `json:"created"`
	Privacy         string      `json:"privacy"`
	FileCount       string      `json:"file_count"`
	FolderCount     string      `json:"folder_count"`
	Revision        string      `json:"revision"`
	DropboxEnabled  string      `json:"dropbox_enabled"`
	Flag            string      `json:"flag"`
	Permissions     interface{} `json:"permissions"`
	UploadKey       string      `json:"upload_key"`
}

func (m *FolderCreate) getCommon() Common {
	return m.Response.Common
}

// UploadPollUpload - структура загрузки Poll
type UploadPollUpload struct {
	Response UploadPollUploadResponse `json:"response"`
}

// UploadPollUploadResponse - структура ответа загрузки Poll
type UploadPollUploadResponse struct {
	Action            string                           `json:"action"`
	Result            string                           `json:"result"`
	CurrentAPIVersion string                           `json:"current_api_version"`
	Doupload          UploadPollUploadResponseDoupload `json:"doupload"`
}

// UploadPollUploadResponseDoupload - структура Doupload
type UploadPollUploadResponseDoupload struct {
	Result      string    `json:"result"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Quickkey    string    `json:"quickkey"`
	Hash        string    `json:"hash"`
	Filename    string    `json:"filename"`
	Size        string    `json:"size"`
	Created     string    `json:"created"`
	Revision    string    `json:"revision"`
	CreatedUtc  time.Time `json:"created_utc"`
}
