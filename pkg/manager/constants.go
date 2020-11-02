package manager

import "time"

// константы
const (
	CookieTokenName                                  = "refreshToken"
	CookiePath                                       = "/api/v1"
	CookieDomain                                     = "localhost"
	LimitDefault                                     = 10
	MinLenPassword                                   = 6
	ImgPathPattern                                   = `^/[0-9a-z]+/[0-9a-z]+\.(jpg|png)$`
	PhonePattern                                     = `^(7|9)\d{10,11}$`
	AccessTokenTimeSecond              time.Duration = 60 * 10           // 10 минут
	RefreshTokenTimeSecond             time.Duration = 60 * 60 * 24 * 30 // месяц
	SessionLimit                                     = 4                 // один прибавляется
	IsAdmin                                          = "admin"
	DirImages                                        = "./web/images"
	DirResample                                      = "./web/resample"
	DirEmail                                         = "./web/email"
	Domain                                           = "https://www.altair.uz"
	HashLen                                          = 32
	MinSecBetweenSendSmsForVerifyPhone               = 60
	MinSecLifeVerifyCode                             = 60 * 20 // 20 минут
)

// TagKindNumber - теги, в которых значения являются числами
var TagKindNumber = []string{"checkbox", "radio", "select", "input_number", "photo"}

// AvailableKindSoc - разрешенные метки для соц. сетей
var AvailableKindSoc = []string{"vk", "ok", "fb", "ggl"}
