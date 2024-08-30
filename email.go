package main

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

func sendEmail(msg string) {
	addr := emne.Config.MailHost + ":" + emne.Config.MailPort
	auth := smtp.PlainAuth("", emne.Config.MailFrom, emne.Config.MailPass, emne.Config.MailHost)
	to := []string{emne.Config.MailTo}

	// Making sure the emails are rendered as HTML.
	msg = fmt.Sprintf("To: %s\nFrom: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=UTF-8;\n", emne.Config.MailTo, emne.Config.MailFrom, "News") + msg

	err := smtp.SendMail(addr, auth, emne.Config.MailFrom, to, []byte(msg))
	if err != nil {
		slog.Warn("Sending mail", "err", err.Error())
		return
	}
}
