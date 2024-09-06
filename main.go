package main

import (
	"email-news/apis"
	"fmt"
	"log/slog"
	"net/http"
)

var emne apis.EmailNews

func main() {
	var err error

	emne, err = apis.Init()
	if err != nil {
		panic("Could not initialize.")
	}

	go scheduleScrape()

	router := http.NewServeMux()

	router.HandleFunc("/", serveHome)
	router.HandleFunc("/{year}/{month}/{day}", serveDateNews)

	slog.Info("Starting server", "port", emne.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", emne.Config.Port), router); err != nil {
		slog.Error("Error while server listening: " + err.Error())
		return
	}
}
