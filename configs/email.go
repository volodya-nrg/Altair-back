package configs

// EmailConfig - структура конфига для е-мэйла (SMTP)
type EmailConfig struct {
	From       string `json:"from"`
	FromName   string `json:"fromName"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Port       int    `json:"port"` // именно int
	SMTPServer string `json:"smtpServer"`
}
