package email

import (
	"crypto/tls"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
}

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{
		SMTPInfo: info,
	}
}

func NewDefaultEmail() *Email {
	container := global.Engine.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	return &Email{
		SMTPInfo: &SMTPInfo{
			Host:     configService.GetString("email.host"),
			Port:     configService.GetInt("email.port"),
			IsSSL:    configService.GetBool("email.isssl"),
			UserName: configService.GetString("email.username"),
			Password: configService.GetString("email.password"),
			From:     configService.GetString("email.from"),
		},
	}
}

func (e *Email) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
