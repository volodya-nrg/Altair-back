package configs

// MediafireConfig - структура конфига для стораджа Mediafire
type MediafireConfig struct {
	AppID        string `json:"appID"`
	AppName      string `json:"appName"`
	APIKey       string `json:"apiKey"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	Domain       string `json:"domain"`
	FolderKey    string `json:"folderKey"`
}
