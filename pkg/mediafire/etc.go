package mediafire

type MediafireCommonInterface interface {
	getMediafireCommon() MediafireCommon
}

type MediafireCommon struct {
	Error             int    `json:"error"`
	Message           string `json:"message"`
	Deprecated        string `json:"deprecated"`
	Result            string `json:"result"`
	NewKey            string `json:"new_key"`
	CurrentApiVersion string `json:"current_api_version"`
	Action            string `json:"action"`
}

type MediafireSessionToken struct {
	Response MediafireSessionTokenResponse `json:"response"`
}
type MediafireSessionTokenResponse struct {
	MediafireCommon
	SessionToken string `json:"session_token"`
	SecretKey    string `json:"secret_key"`
	SecretKeyInt uint64
	Time         string `json:"time"`
	Ekey         string `json:"ekey"`
	Pkey         string `json:"pkey"`
}

func (m MediafireSessionToken) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadGetOptions struct {
	Response MediafireUploadGetOptionsResponse `json:"response"`
}
type MediafireUploadGetOptionsResponse struct {
	MediafireCommon
	DisableFlash         string `json:"disable_flash"`
	DisableHtml5         string `json:"disable_html5"`
	DisableInstant       string `json:"disable_instant"`
	ActionOnDuplicate    string `json:"action_on_duplicate"`
	UsedStorageSize      int    `json:"used_storage_size"`
	StorageLimit         int    `json:"storage_limit"`
	StorageLimitExceeded string `json:"storage_limit_exceeded"`
}

func (m MediafireUploadGetOptions) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadSetOptions struct {
	Response MediafireUploadSetOptionsResponse `json:"response"`
}
type MediafireUploadSetOptionsResponse struct {
	MediafireCommon
}

func (m MediafireUploadSetOptions) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadAddWebUpload struct {
	Response MediafireUploadAddWebUploadResponse `json:"response"`
}
type MediafireUploadAddWebUploadResponse struct {
	MediafireCommon
	UploadKey string `json:"upload_key"`
}

func (m MediafireUploadAddWebUpload) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUserGetInfo struct {
	Response MediafireUserGetInfoResponse `json:"response"`
}
type MediafireUserGetInfoResponse struct {
	MediafireCommon
	UserInfo map[string]interface{} `json:"user_info"`
}

func (m MediafireUserGetInfo) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadSimple struct {
	Response MediafireUploadSimpleResponse `xml:"response"`
}
type MediafireUploadSimpleResponse struct {
	MediafireCommon
	Doupload          MediafireUploadSimpleResponseDoupload `xml:"doupload"`
	Server            string                                `xml:"server"`
	NewDeviceRevision int                                   `xml:"new_device_revision"`
}
type MediafireUploadSimpleResponseDoupload struct {
	Result int    `xml:"result"`
	Key    string `xml:"key"`
}

func (m MediafireUploadSimple) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadCheck struct {
	Response MediafireUploadCheckResponse `xml:"response"`
}
type MediafireUploadCheckResponse struct {
	MediafireCommon
	AvailableSpace       string                                `json:"available_space"`
	FileExists           string                                `json:"file_exists"`
	HashExists           string                                `json:"hash_exists"`
	StorageLimit         string                                `json:"storage_limit"`
	StorageLimitExceeded string                                `json:"storage_limit_exceeded"`
	UnitSize             string                                `json:"unit_size"`
	UploadUrl            MediafireUploadCheckResponseUploadUrl `json:"upload_url"`
	UsedStorageSize      string                                `json:"used_storage_size"`
}
type MediafireUploadCheckResponseUploadUrl struct {
	Resumable string `json:"resumable"`
	Simple    string `json:"simple"`
}

func (m MediafireUploadCheck) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUploadInstant struct {
	Response MediafireUploadInstantResponse `xml:"response"`
}
type MediafireUploadInstantResponse struct {
	MediafireCommon
	Quickkey string `json:"quickkey"`
	Filename string `json:"filename"`
}

func (m MediafireUploadInstant) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUserGetActionToken struct {
	Response MediafireUserGetActionTokenResponse `json:"response"`
}
type MediafireUserGetActionTokenResponse struct {
	MediafireCommon
	ActionToken string `json:"action_token"`
}

func (m MediafireUserGetActionToken) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUserDestroyActionToken struct {
	Response MediafireUserDestroyActionTokenResponse `json:"response"`
}
type MediafireUserDestroyActionTokenResponse struct {
	MediafireCommon
}

func (m MediafireUserDestroyActionToken) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireUserRenewSessionToken struct {
	Response MediafireUserRenewSessionTokenResponse `json:"response"`
}
type MediafireUserRenewSessionTokenResponse struct {
	MediafireCommon
	SessionToken string `json:"session_token"`
}

func (m MediafireUserRenewSessionToken) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireFolderGetInfo struct {
	Response MediafireFolderGetInfoResponse `json:"response"`
}
type MediafireFolderGetInfoResponse struct {
	MediafireCommon
	FolderInfo map[string]string `json:"folder_info"`
}

func (m MediafireFolderGetInfo) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}

//----------------------------------------------------------------------------------------------------------------------
type MediafireFolderConfigureFiledrop struct {
	Response MediafireFolderConfigureFiledropResponse `json:"response"`
}
type MediafireFolderConfigureFiledropResponse struct {
	MediafireCommon
	FiledropKey       string `json:"filedrop_key"`
	HostedFileDrop    string `json:"hosted_file_drop"` // "http://www.mediafire.com/filedrop/filedrop_hosted.php?drop=9ed8e61f15b6101b486df3cfa8a31e77e5cdc216c86fbf96d1d75e124ecc9dfd",
	HtmlEmbedCode     string `json:"html_embed_code"`
	NewDeviceRevision int    `json:"new_device_revision"`
}

func (m MediafireFolderConfigureFiledrop) getMediafireCommon() MediafireCommon {
	return m.Response.MediafireCommon
}
