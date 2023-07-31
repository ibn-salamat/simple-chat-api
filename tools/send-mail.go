package tools

import (
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"net/smtp"
)

func SendMail(to string, subject string, content string) error {
	message := []byte(fmt.Sprintf("Subject: %s \r\n\r\n\n %s", subject, content))
	addr := fmt.Sprintf("%s:%s", config.EnvData.SMTP_ADDR, ":587")
	auth := smtp.PlainAuth(
		"",
		config.EnvData.SMTP_LOGIN,
		config.EnvData.SMTP_KEY,
		config.EnvData.SMTP_ADDR,
	)
	from := config.EnvData.SMTP_LOGIN

	return smtp.SendMail(
		addr,
		auth,
		from,
		[]string{to},
		message,
	)

}
