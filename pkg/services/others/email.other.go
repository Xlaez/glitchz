package others

import (
	"glitchz/pkg/utils"
	"net/smtp"
)

func SendEmail(config utils.Config, to []string, message []byte) error {
	username := config.SmtpUsername
	password := config.SmtpPassword
	smtpHost := config.SmtpHost

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := config.AppEmail
	smtpUrl := smtpHost + ":25"

	if err := smtp.SendMail(smtpUrl, auth, from, to, message); err != nil {
		return err
	}

	return nil
}
