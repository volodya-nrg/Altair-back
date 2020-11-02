package emailer

import (
	"altair/configs"
	"altair/pkg/manager"
	"bytes"
	"github.com/go-gomail/gomail"
	"html/template"
)

// NewEmailRequest - фабрика, создание объекта е-мэйла
func NewEmailRequest(to, subject, body string) *EmailRequest {
	result := new(EmailRequest)

	result.to = to
	result.subject = subject
	result.body = body

	return result
}

// EmailRequest - структура для отправки е-мэйла
type EmailRequest struct {
	to      string
	subject string
	body    string
}

// SendMail - метод отправки е-мэйла
func (er *EmailRequest) SendMail() (bool, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", configs.Cfg.Email.From)
	m.SetHeader("To", er.to)
	m.SetHeader("Subject", er.subject)
	m.SetBody("text/html", er.body)
	d := gomail.NewDialer(
		configs.Cfg.Email.SMTPServer,
		configs.Cfg.Email.Port,
		configs.Cfg.Email.Login,
		configs.Cfg.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}

// ParseTemplate - ф-ия парсинга шаблона и вставки в нее данных
func (er *EmailRequest) ParseTemplate(tplFileName string, data interface{}) error {
	path := manager.DirEmail + "/templates" // именно web, т.к. она открыта для Докера
	files := []string{path + "/layouts/layout-1.html", path + "/pages/" + tplFileName}

	for _, v := range files {
		if !manager.FolderOrFileExists(v) {
			return manager.ErrFileNotFound
		}
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, data); err != nil {
		return err
	}

	er.body = buf.String()

	return nil
}
