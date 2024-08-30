package main

import (
	"bytes"
	"context"
	"log/slog"
	"text/template"
	"time"

	"email-news/server"

	"github.com/a-h/templ"
	"github.com/go-co-op/gocron"
)

func scheduleScrape() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Days().At(emne.Config.Time).Do(scrape)

	s.StartAsync()
}

func scrape() {
	var n server.News

	dateString := server.GetDateNowString()

	n, err := server.GetHNContent(n)
	if err != nil {
		slog.Warn("Getting hacker news content", "error", err.Error())
	}

	n, err = server.GetTLDRContent(n, server.Tech, dateString)
	if err != nil {
		slog.Warn("Getting tldr tech content", "error", err.Error())
	}

	n, err = server.GetTLDRContent(n, server.WebDev, dateString)
	if err != nil {
		slog.Warn("Getting tldr webdev content", "error", err.Error())
	}

	n, err = server.GetMBContent(n)
	if err != nil {
		slog.Warn("Getting morning brew content", "error", err.Error())
	}

	err = emne.StoreNews(n)
	if err != nil {
		slog.Warn("Storing news", "err", err.Error())
	}

	processEmail(n, dateString[0:3], dateString[5:6], dateString[8:9])
}

func processEmail(n server.News, year, month, day string) {
	parser := template.Must(template.New("").Parse(`{{ . }}`))

	formattedMB, err := formatMB(n)
	if err != nil {
		slog.Warn("Formatting MB data", "err", err.Error())
		return
	}

	component := NewsEmail(n, formattedMB, year, month, day) // TODO - actually use YMD in the template.

	html, err := templ.ToGoHTML(context.Background(), component)
	if err != nil {
		slog.Error("Generating html", "err", err.Error())
		return
	}

	buf := new(bytes.Buffer)
	err = parser.Execute(buf, html)
	if err != nil {
		slog.Error("Executing template", "err", err.Error())
		return
	}

	message := buf.String()

	sendEmail(message)
}
