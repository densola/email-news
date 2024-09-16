package main

import (
	"bytes"
	"context"
	"log/slog"
	"strconv"
	"text/template"
	"time"

	"email-news/apis"

	"github.com/a-h/templ"
	"github.com/go-co-op/gocron"
)

func scheduleScrape() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Days().At(emne.Config.Time).Do(scrape)

	s.StartAsync()
}

func scrape() {
	var n apis.News

	date := apis.GetDateNowString()
	dbDate, err := strconv.Atoi(date)
	if err != nil {
		slog.Warn("Converting date string to int", "err", err.Error())
	}

	n, err = apis.GetHNContent(n)
	if err != nil {
		slog.Warn("Getting hacker news content", "error", err.Error())
	}

	n, err = apis.GetTLDRContent(n, apis.Tech, date)
	if err != nil {
		slog.Warn("Getting tldr tech content", "error", err.Error())
	}

	n, err = apis.GetTLDRContent(n, apis.WebDev, date)
	if err != nil {
		slog.Warn("Getting tldr webdev content", "error", err.Error())
	}

	n, err = apis.GetMBContent(n)
	if err != nil {
		slog.Warn("Getting morning brew content", "error", err.Error())
	}

	err = emne.StoreNews(n, dbDate)
	if err != nil {
		slog.Warn("Storing news", "err", err.Error())
	}

	parser := template.Must(template.New("").Parse(`{{ . }}`))

	formattedMB, err := formatMB(n)
	if err != nil {
		slog.Warn("Formatting MB data", "err", err.Error())
		return
	}

	component := NewsEmail(n, formattedMB, date[0:4], date[4:6], date[6:8]) // TODO - actually use YMD in the template.

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

	emne.SendEmail(message)
}
