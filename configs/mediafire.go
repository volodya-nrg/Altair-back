package configs

type MediafireConfig struct {
	AppID        string `json:"appID"`
	AppName      string `json:"appName"`
	APIKey       string `json:"APIKey"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	Domain       string `json:"domain"`
	FolderKey    string `json:"folderKey"`
}
