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

func (m SessionToken) getCommon() Common {
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

func (m UploadSimple) getCommon() Common {
	return m.Response.Common
}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUploadGetOptions struct {
//	Response MediafireUploadGetOptionsResponse `json:"response"`
//}
//type MediafireUploadGetOptionsResponse struct {
//	MediafireCommon
//	DisableFlash         string `json:"disable_flash"`
//	DisableHtml5         string `json:"disable_html5"`
//	DisableInstant       string `json:"disable_instant"`
//	ActionOnDuplicate    string `json:"action_on_duplicate"`
//	UsedStorageSize      int    `json:"used_storage_size"`
//	StorageLimit         int    `json:"storage_limit"`
//	StorageLimitExceeded string `json:"storage_limit_exceeded"`
//}
//
//func (m MediafireUploadGetOptions) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUploadSetOptions struct {
//	Response MediafireUploadSetOptionsResponse `json:"response"`
//}
//type MediafireUploadSetOptionsResponse struct {
//	MediafireCommon
//}
//
//func (m MediafireUploadSetOptions) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUserGetInfo struct {
//	Response MediafireUserGetInfoResponse `json:"response"`
//}
//type MediafireUserGetInfoResponse struct {
//	MediafireCommon
//	UserInfo map[string]interface{} `json:"user_info"`
//}
//
//func (m MediafireUserGetInfo) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUploadCheck struct {
//	Response MediafireUploadCheckResponse `xml:"response" json:"response"`
//}
//type MediafireUploadCheckResponse struct {
//	MediafireCommon
//	AvailableSpace       string                                `json:"available_space"`
//	FileExists           string                                `json:"file_exists"`
//	HashExists           string                                `json:"hash_exists"`
//	StorageLimit         string                                `json:"storage_limit"`
//	StorageLimitExceeded string                                `json:"storage_limit_exceeded"`
//	UnitSize             string                                `json:"unit_size"`
//	UploadUrl            MediafireUploadCheckResponseUploadUrl `json:"upload_url"`
//	UsedStorageSize      string                                `json:"used_storage_size"`
//}
//type MediafireUploadCheckResponseUploadUrl struct {
//	Resumable string `json:"resumable"`
//	Simple    string `json:"simple"`
//}
//
//func (m MediafireUploadCheck) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUserGetActionToken struct {
//	Response MediafireUserGetActionTokenResponse `json:"response"`
//}
//type MediafireUserGetActionTokenResponse struct {
//	MediafireCommon
//	ActionToken string `json:"action_token"`
//}
//
//func (m MediafireUserGetActionToken) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUserDestroyActionToken struct {
//	Response MediafireUserDestroyActionTokenResponse `json:"response"`
//}
//type MediafireUserDestroyActionTokenResponse struct {
//	MediafireCommon
//}
//
//func (m MediafireUserDestroyActionToken) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUserRenewSessionToken struct {
//	Response MediafireUserRenewSessionTokenResponse `json:"response"`
//}
//type MediafireUserRenewSessionTokenResponse struct {
//	MediafireCommon
//	SessionToken string `json:"session_token"`
//}
//
//func (m MediafireUserRenewSessionToken) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//----------------------------------------------------------------------------------------------------------------------
//type MediafireFolderGetInfo struct {
//	Response MediafireFolderGetInfoResponse `json:"response"`
//}
//type MediafireFolderGetInfoResponse struct {
//	MediafireCommon
//	FolderInfo map[string]string `json:"folder_info"`
//}
//
//func (m MediafireFolderGetInfo) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

//type MediafireFolderConfigureFiledrop struct {
//	Response MediafireFolderConfigureFiledropResponse `json:"response"`
//}
//type MediafireFolderConfigureFiledropResponse struct {
//	MediafireCommon
//	FiledropKey       string `json:"filedrop_key"`
//	HostedFileDrop    string `json:"hosted_file_drop"` // "http://www.mediafire.com/filedrop/filedrop_hosted.php?drop=9ed8e61f15b6101b486df3cfa8a31e77e5cdc216c86fbf96d1d75e124ecc9dfd",
//	HtmlEmbedCode     string `json:"html_embed_code"`
//	NewDeviceRevision int    `json:"new_device_revision"`
//}
//
//func (m MediafireFolderConfigureFiledrop) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}

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

func (m FolderCreate) getCommon() Common {
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

//----------------------------------------------------------------------------------------------------------------------
//type MediafireUploadResumable struct {
//	Response MediafireUploadResumableResponse `xml:"response" json:"response"`
//}
//type MediafireUploadResumableResponse struct {
//	MediafireCommon
//	Doupload          MediafireUploadResumableResponseDoupload `json:"doupload"`
//	Server            string                                   `json:"server"`
//	ResumableUpload   interface{}                              `json:"resumable_upload"`
//	NewDeviceRevision int                                      `json:"new_device_revision"`
//}
//type MediafireUploadResumableResponseDoupload struct {
//	Result string `json:"result"`
//	Key    string `json:"key"`
//}
//
//func (m MediafireUploadResumable) getMediafireCommon() MediafireCommon {
//	return m.Response.MediafireCommon
//}
