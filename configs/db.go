package configs

type DbConfig struct {
	Host     string `json:"Host"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Port     uint   `json:"Port"`
	User     string `json:"User"`
}
