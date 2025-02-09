package email

import (
	"crypto/tls"
	"diplom_chat_ssh/internal/model"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

func SendMail(receiver string, subject string, message string, configEmail model.Email) {
	d := gomail.NewDialer(
		configEmail.Host,
		configEmail.Port,
		configEmail.Username,
		configEmail.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", configEmail.Username)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	if err := d.DialAndSend(m); err != nil {
		log.Error().Err(err).Msg("")
	}
}
