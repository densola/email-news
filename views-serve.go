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
		slog.Warn("Getting links for display on homepage", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = HomePage(links).Render(r.Context(), w)
	if err != nil {
		slog.Warn("Rendering home page", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serveDateNews(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("year")
	month := r.PathValue("month")
	day := r.PathValue("day")
	ymdDate := year + "/" + month + "/" + day

	news, err := emne.GetNews(year, month, day)
	if err != nil {
		slog.Warn("Getting news based on date", "err", err, "ymd date", ymdDate)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	formattedMB, err := formatMB(news)
	if err != nil {
		slog.Warn("Preserving line breaks for morning brew paragraphs", "err", err, "ymd date", ymdDate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := strconv.Atoi(month)
	if err != nil {
		slog.Warn("Getting int from month string", "err", err, "month", month)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = NewsOnDatePage(news, formattedMB, year, time.Month(m).String(), day).Render(r.Context(), w)
	if err != nil {
		slog.Warn("Rendering news page", "err", err, "ymd date", ymdDate)
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

		// Create a component containing raw HTML
		content := raw(buf.String())
		fmb = append(fmb, content)
	}

	return fmb, nil
}

func raw(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
