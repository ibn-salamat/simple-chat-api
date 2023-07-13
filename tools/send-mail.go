package tools

import (
	"fmt"
	"ibn-salamat/simple-chat-api/config"
	"net/smtp"
)

func SendMail(to string, subject string, content string) error {
	message := []byte(fmt.Sprintf("Subject: %s \r\n\r\n\n %s", subject, content))
	addr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth(
		"",
		"n.salamatoff@gmail.com",
		config.EnvData.GOOGLE_GMAIL_KEY,
		"smtp.gmail.com",
	)
	from := "admin@simple-chat.com"

	return smtp.SendMail(
		addr,
		auth,
		from,
		[]string{to},
		message,
	)

}
