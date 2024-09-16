package apis

import (
	"fmt"
	"net/smtp"
)

func (e EmailNews) SendEmail(msg string) error {
	addr := e.Config.MailHost + ":" + e.Config.MailPort
	auth := smtp.PlainAuth("", e.Config.MailFrom, e.Config.MailPass, e.Config.MailHost)
	to := []string{e.Config.MailTo}

	// Making sure the emails are rendered as HTML.
	msg = fmt.Sprintf("To: %s\nFrom: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=UTF-8;\n", e.Config.MailTo, e.Config.MailFrom, "News") + msg

	err := smtp.SendMail(addr, auth, e.Config.MailFrom, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
