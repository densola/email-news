package main

import (
	"bytes"
	"context"
	"email-news/apis"
	"html/template"
	"log/slog"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-co-op/gocron"
)

func scheduleDailyNews() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Days().At(emne.Config.Time).Do(handleDailyNews)

	s.StartAsync()
}

func handleDailyNews() {
	var n apis.News

	date := apis.GetDateNowString()

	n, err := apis.GetContent(emne.Config.WeatherAPIKey, emne.Config.WeatherAPILocation, date)
	if err != nil {
		slog.Warn("Getting content for date", "err", err, "date", date)
		return
	}

	dbDate, err := strconv.Atoi(date)
	if err != nil {
		slog.Warn("Converting date string to int", "err", err, "date", date)
		return
	}

	err = emne.StoreNews(n, dbDate)
	if err != nil {
		slog.Warn("Storing news", "err", err)
		return
	}

	parser := template.Must(template.New("").Parse(`{{ . }}`))

	formattedMB, err := formatMB(n)
	if err != nil {
		slog.Warn("Preserving mb line breaks for emails", "err", err)
		return
	}

	month, err := parseMonth(date[4:6])
	if err != nil {
		slog.Warn("Getting name for month", "err", err, "month", date[4:6])
	}

	component := NewsEmail(n, formattedMB, date[0:4], month, date[6:8])

	html, err := templ.ToGoHTML(context.Background(), component)
	if err != nil {
		slog.Warn("Generating html from templ component", "err", err)
		return
	}

	buf := new(bytes.Buffer)
	err = parser.Execute(buf, html)
	if err != nil {
		slog.Warn("Applying template to html", "err", err)
		return
	}

	message := buf.String()

	emne.SendEmail(message)
}
