package service

import (
	"altair/pkg/manager"
	"fmt"
)

// NewSMSService - фабрика, создает объект СМС
func NewSMSService(apiKeySrc string, domain string) *SMSService {
	sms := &SMSService{
		apiKey: apiKeySrc,
		domain: domain,
	}

	return sms
}

// SMSService - структура СМС
type SMSService struct {
	apiKey string
	domain string
}

// Send - отправить данные на удаленный сервис
func (sms SMSService) Send(data SMSSendRequest) (*SMSSendResponse, error) {
	receiver := new(SMSSendResponse)
	query := map[string]string{
		"to":     data.To,
		"msg":    data.Msg,
		"from":   data.From,
		"test":   fmt.Sprintf("%d", data.Test),
		"api_id": sms.apiKey,
		"json":   "1",
	}

	if err := manager.MakeRequest("post", sms.domain+"/sms/send", receiver, query); err != nil {
		return receiver, err
	}

	return receiver, nil
}

// Balance - узнать баланс
func (sms SMSService) Balance() (*SMSBalanceResponse, error) {
	receiver := new(SMSBalanceResponse)
	query := map[string]string{
		"api_id": sms.apiKey,
		"json":   "1",
	}

	if err := manager.MakeRequest("get", sms.domain+"/my/balance", receiver, query); err != nil {
		return receiver, err
	}

	return receiver, nil
}

// Cost - цена за смс
func (sms SMSService) Cost(data SMSCostRequest) (*SMSCostResponse, error) {
	receiver := new(SMSCostResponse)
	query := map[string]string{
		"to":     data.To,
		"msg":    data.Msg,
		"api_id": sms.apiKey,
		"json":   "1",
	}

	if err := manager.MakeRequest("post", sms.domain+"/sms/cost", receiver, query); err != nil {
		return receiver, err
	}

	return receiver, nil
}

// GetStatusInfo - получение описания статуса
func (sms SMSService) GetStatusInfo(statusCode int) string {
	var result string

	switch statusCode {
	case -1:
		result = "Сообщение не найдено"
	case 100:
		result = "Запрос выполнен или сообщение находится в нашей очереди"
	case 101:
		result = "Сообщение передается оператору"
	case 102:
		result = "Сообщение отправлено (в пути)"
	case 103:
		result = "Сообщение доставлено"
	case 104:
		result = "Не может быть доставлено: время жизни истекло"
	case 105:
		result = "Не может быть доставлено: удалено оператором"
	case 106:
		result = "Не может быть доставлено: сбой в телефоне"
	case 107:
		result = "Не может быть доставлено: неизвестная причина"
	case 108:
		result = "Не может быть доставлено: отклонено"
	case 110:
		result = "Сообщение прочитано (для Viber, временно не работает)"
	case 150:
		result = "Не может быть доставлено: не найден маршрут на данный номер"
	case 200:
		result = "Неправильный api_id"
	case 201:
		result = "Не хватает средств на лицевом счету"
	case 202:
		result = "Неправильно указан номер телефона получателя, либо на него нет маршрута"
	case 203:
		result = "Нет текста сообщения"
	case 204:
		result = "Имя отправителя не согласовано с администрацией"
	case 205:
		result = "Сообщение слишком длинное (превышает 8 СМС)"
	case 206:
		result = "Будет превышен или уже превышен дневной лимит на отправку сообщений"
	case 207:
		result = "На этот номер нет маршрута для доставки сообщений"
	case 208:
		result = "Параметр time указан неправильно"
	case 209:
		result = "Вы добавили этот номер (или один из номеров) в стоп-лист"
	case 210:
		result = "Используется GET, где необходимо использовать POST"
	case 211:
		result = "Метод не найден"
	case 212:
		result = "Текст сообщения необходимо передать в кодировке UTF-8 (вы передали в другой кодировке)"
	case 213:
		result = "Указано более 100 номеров в списке получателей"
	case 220:
		result = "Сервис временно недоступен, попробуйте чуть позже"
	case 230:
		result = "Превышен общий лимит количества сообщений на этот номер в день"
	case 231:
		result = "Превышен лимит одинаковых сообщений на этот номер в минуту"
	case 232:
		result = "Превышен лимит одинаковых сообщений на этот номер в день"
	case 233:
		result = "Превышен лимит отправки повторных сообщений с кодом на этот номер за короткий промежуток времени (\"защита от мошенников\", можно отключить в разделе \"Настройки\")"
	case 300:
		result = "Неправильный token (возможно истек срок действия, либо ваш IP изменился)"
	case 301:
		result = "Неправильный api_id, либо логин/пароль"
	case 302:
		result = "Пользователь авторизован, но аккаунт не подтвержден (пользователь не ввел код, присланный в регистрационной смс)"
	case 303:
		result = "Код подтверждения неверен"
	case 304:
		result = "Отправлено слишком много кодов подтверждения. Пожалуйста, повторите запрос позднее"
	case 305:
		result = "Слишком много неверных вводов кода, повторите попытку позднее"
	case 500:
		result = "Ошибка на сервере. Повторите запрос."
	case 901:
		result = "Callback: URL неверный (не начинается на http://)"
	case 902:
		result = "Callback: Обработчик не найден (возможно был удален ранее)"
	default:
		result = "Не известный статус кода"
	}

	return result
}

// SMSSendRequest - структура запроса отправки
type SMSSendRequest struct {
	To   string `json:"to"`
	Msg  string `json:"msg"`
	From string `json:"from"`
	Test int    `json:"test"`
}

// SMSSendResponse - структура приема
type SMSSendResponse struct {
	Balance    float64                       `json:"balance"`
	Sms        map[string]SMSSendResponseSms `json:"sms"`
	Status     string                        `json:"status"`
	StatusCode int                           `json:"status_code"`
}

// SMSSendResponseSms - структура приема СМС
type SMSSendResponseSms struct {
	Cost       string `json:"cost"`
	SmsID      string `json:"sms_id"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
}

// SMSBalanceResponse - структура приема баланса
type SMSBalanceResponse struct {
	Balance    float64 `json:"balance"`
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
}

// SMSCostRequest - структура запроса отправки на цену
type SMSCostRequest struct {
	To  string `json:"to"`
	Msg string `json:"msg"`
}

// SMSCostResponse - структура приема данных о цене
type SMSCostResponse struct {
	Status     string                        `json:"status"`
	StatusCode int                           `json:"status_code"`
	Sms        map[string]SMSCostResponseSms `json:"sms"`
	TotalCost  float64                       `json:"total_cost"`
	TotalSms   int                           `json:"total_sms"`
}

// SMSCostResponseSms - структура приема об СМС
type SMSCostResponseSms struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Cost       float64 `json:"cost"`
	Sms        int     `json:"sms"`
	StatusText string  `json:"status_text"`
}

// private -------------------------------------------------------------------------------------------------------------
