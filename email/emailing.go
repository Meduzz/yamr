package email

import (
	"os"
	"net/smtp"
	"fmt"
)

type (
	Email struct {
		To string
		Subject string
		Body string
		From string
	}
)

// TODO move these to a generic settings store, administered from a web ui.
var ident = fromEnv("MAIL_IDENT", "")
var username = fromEnv("MAIL_USERNAME", "")
var password = fromEnv("MAIL_PASSWORD", "")
var mail_server = fromEnv("MAIL_SERVER", "")
var mail_port = fromEnv("MAIL_PORT", "")

var auth = smtp.PlainAuth(ident, username, password, mail_server)

func SendMail(mail *Email) error {
	return  smtp.SendMail(fmt.Sprintf("%s:%s", mail_server, mail_port), auth, mail.From, mail.To, mail.ToMsg())
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}

func (e *Email) ToMsg() []byte {
	return []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", e.To, e.Subject, e.Body))
}
