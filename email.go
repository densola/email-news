package main

import (
	"log/slog"
	"net/smtp"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func email(msg string) {
	from := emne.Config.MailFrom
	password := emne.Config.MailPass

	to := []string{emne.Config.MailTo}

	smtpServer := smtpServer{host: emne.Config.MailHost, port: emne.Config.MailPort}

	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	err := smtp.SendMail(smtpServer.Address(), auth, from, to, []byte(msg))
	if err != nil {
		slog.Warn("Sending mail", "err", err.Error())
		return
	}
}
