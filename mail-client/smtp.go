package mailclient

import (
	"fmt"
	"net/smtp"
)

func SendMail(from string, to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	var auth smtp.Auth
	if config.Username != "" && config.Password != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	msg := buildMessage(from, to, subject, body)

	return smtp.SendMail(addr, auth, from, to, []byte(msg))
}
