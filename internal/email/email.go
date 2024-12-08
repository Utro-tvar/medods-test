package email

import "net/smtp"

func Send(to string, text []byte) error {
	from := "service@gmail.com"
	password := "<Email Password>"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, text)
	return nil
}
