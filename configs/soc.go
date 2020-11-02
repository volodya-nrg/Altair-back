package configs

// SocialsConfig - общая структура конфигов соц. сетей
type SocialsConfig struct {
	Vk  SocVkConfig  `json:"vk"`
	Ok  SocOkConfig  `json:"ok"`
	Fb  SocFbConfig  `json:"fb"`
	Ggl SocGglConfig `json:"ggl"`
}

// SocVkConfig - структура конфига VK
type SocVkConfig struct {
	ClientID     uint64 `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

// SocOkConfig - структура конфига OK
type SocOkConfig struct {
	ClientID     uint64 `json:"clientID"`
	ClientPublic string `json:"clientPublic"`
	ClientSecret string `json:"clientSecret"`
}

// SocFbConfig - структура конфига FB
type SocFbConfig struct {
	ClientID     uint64 `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

// SocGglConfig - структура конфига GGL
type SocGglConfig struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}
