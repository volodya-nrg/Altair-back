package configs

// DbConfig - структура конфига для DB
type DbConfig struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
}
