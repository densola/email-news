package main

import (
	"bytes"
	"context"
	"email-news/apis"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	links, err := emne.GetHomeLinks()
	if err != nil {
		slog.Error("Getting dates with news recorded", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = HomePage(links).Render(r.Context(), w)
	if err != nil {
		slog.Error("Rendering home page", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serveDateNews(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("year")
	month := r.PathValue("month")
	day := r.PathValue("day")

	m, err := strconv.Atoi(month)
	if err != nil {
		slog.Warn("Converting month string to int", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	news, err := emne.GetNews(year, month, day)
	if err != nil {
		slog.Error("Getting news based on date", "err", err.Error(), "year", year, "month", month, "day", day)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	formattedMB, err := formatMB(news)
	if err != nil {
		slog.Warn("Formatting MB data", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = NewsOnDatePage(news, formattedMB, year, time.Month(m).String(), day).Render(r.Context(), w)
	if err != nil {
		slog.Error("Rendering news page", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func formatMB(n apis.News) ([]templ.Component, error) {
	fmb := []templ.Component{}

	for i := 0; i < len(n.MBArticles); i++ {
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(n.MBArticles[i].Overview), &buf); err != nil {
			return nil, fmt.Errorf("converting markdown to html: %w", err)
		}

		// Create a component containing raw HTML.
		content := Raw(buf.String())
		fmb = append(fmb, content)
	}

	return fmb, nil
}

func Raw(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
