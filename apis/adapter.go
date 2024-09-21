/*
adapter.go handles calls from the view before passing them into the model.
*/
package apis

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type EmailNews struct {
	model  Model
	Config config
}

type config struct {
	Port     int    `env:"PORT"`
	Time     string `env:"TIME"`
	MailFrom string `env:"MAILFROM"`
	MailTo   string `env:"MAILTO"`
	MailPass string `env:"MAILPASS"`
	MailHost string `env:"MAILHOST"`
	MailPort string `env:"MAILPORT"`
}

type Hyperlink struct {
	Text        string
	Destination string
}

func Initialize() (EmailNews, error) {
	emne := EmailNews{}
	cfg := config{}

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return emne, fmt.Errorf("loading config file: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return emne, fmt.Errorf("parsing config file: %w", err)
	}

	emne.Config = cfg

	db, err := emne.sqliteConnect()
	if err != nil {
		return emne, fmt.Errorf("connecting to sqlite3 db: %w", err)
	}

	emne.model = newModel(db)

	return emne, nil
}

func (na *EmailNews) sqliteConnect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "db.db")
	if err != nil {
		return db, fmt.Errorf("opening connection to sqlite3 db: %w", err)
	}

	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

func (na EmailNews) StoreNews(news News, date int) error {
	b, err := json.Marshal(news)
	if err != nil {
		return fmt.Errorf("encoding news to json: %w", err)
	}

	err = na.model.insertNews(b, date)
	if err != nil {
		return fmt.Errorf("inserting byte and date: %w", err)
	}

	return nil
}

func (na EmailNews) GetNews(year, month, day string) (News, error) {
	var n News

	date := fmt.Sprintf("%s%s%s", year, month, day)
	dbDate, err := strconv.Atoi(date)
	if err != nil {
		return n, fmt.Errorf("converting date string to int: %w", err)
	}

	b, err := na.model.getNews(dbDate)
	if err != nil {
		return n, fmt.Errorf("getting byte based on date: %w", err)
	}

	err = json.Unmarshal(b, &n)
	if err != nil {
		return n, fmt.Errorf("parsing encoded json news: %w", err)
	}

	return n, nil
}

// TODO - If we have a thousand records, should all of them be pulled?
// GetHomeLinks returns a list of all hyperlinks to be presented on the homepage.
func (na EmailNews) GetHomeLinks() ([]Hyperlink, error) {
	var links []Hyperlink
	times, err := na.model.getTimes()
	if err != nil {
		return nil, fmt.Errorf("getting all time values for homepage display: %w", err)
	}

	for _, tm := range times {
		t := fmt.Sprintf("%d", tm)
		y := t[0:4]
		m := t[4:6]
		d := t[6:8]

		var date Hyperlink
		date.Destination = fmt.Sprintf("%s/%s/%s", y, m, d)

		month, err := strconv.Atoi(m)
		if err != nil {
			return nil, fmt.Errorf("converting month string to int: %w", err)
		}

		date.Text = fmt.Sprintf("%s %s, %s", time.Month(month), d, y)

		links = append(links, date)
	}

	return links, nil
}
