package main

import (
	"fmt"
	"log/slog"
	"time"

	"email-news/server"

	"github.com/go-co-op/gocron"
)

func scheduleScrape() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Days().At(emne.Config.Time).Do(scrape)

	s.StartAsync()
}

func scrape() {
	var n server.News

	y, m, d := server.GetYMDNow()

	n, err := server.GetHNContent(n)
	if err != nil {
		slog.Warn("Getting hacker news content", "error", err.Error())
	}

	n, err = server.GetTLDRContent(n, server.Tech, y, m, d)
	if err != nil {
		slog.Warn("Getting tldr tech content", "error", err.Error())
	}

	n, err = server.GetTLDRContent(n, server.WebDev, y, m, d)
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

	emailBody := "Subject: News"

	if len(n.HNThreads.Items) > 0 {
		emailBody += "\n\n\n# HackerNews:\n"
		for _, thread := range n.HNThreads.Items {
			emailBody += fmt.Sprintf("Title: %s\nLink: %s\nComments: %s\n\n", thread.Title, thread.Link, thread.Comments)
		}
	}

	if len(n.TLDRTechArticles) > 0 {
		emailBody += "\n\n\n# TLDR Tech:\n"
		for _, page := range n.TLDRTechArticles {
			emailBody += fmt.Sprintf("Title: %s\nLink: %s\nOverview: %s\n\n", page.Title, page.Link, page.Overview)
		}
	}

	if len(n.TLDRWebDevArticles) > 0 {
		emailBody += "\n\n\n# TLDR WebDev:\n"
		for _, page := range n.TLDRWebDevArticles {
			emailBody += fmt.Sprintf("Title: %s\nLink: %s\nOverview: %s\n\n", page.Title, page.Link, page.Overview)
		}
	}

	if len(n.MBArticles) > 0 {
		emailBody += "\n\n\n# Morning Brew:\n"
		for _, page := range n.MBArticles {
			emailBody += fmt.Sprintf("Title: %s\nOverview: %s\n\n", page.Title, page.Overview)
		}
	}

	email(emailBody)
}
