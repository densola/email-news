package server

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
	News   News
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

func Init() (EmailNews, error) {
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

func (na EmailNews) StoreNews(news News) error {
	b, err := json.Marshal(news)
	if err != nil {
		return fmt.Errorf("marshaling news: %w", err)
	}

	dateNow := GetDateNowString()
	dateNow = dateNow[0:4] + dateNow[5:7] + dateNow[8:10]
	dbDateNow, err := strconv.Atoi(dateNow)
	if err != nil {
		return fmt.Errorf("converting date now string to int: %w", err)
	}

	err = na.model.insertByte(b, dbDateNow)
	if err != nil {
		return fmt.Errorf("inserting news byte and unix now time: %w", err)
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

	b, err := na.model.getByte(dbDate)
	if err != nil {
		return n, fmt.Errorf("getting byte based on date: %w", err)
	}

	err = json.Unmarshal(b, &n)
	if err != nil {
		return n, fmt.Errorf("unmarshaling byte into news variable: %w", err)
	}

	return n, nil
}

// TODO - If we have a thousand records, should all of them be pulled?
// GetHomeLinks returns a list of all hyperlinks to be presented on the homepage.
func (na EmailNews) GetHomeLinks() ([]Hyperlink, error) {
	var links []Hyperlink
	times, err := na.model.getTimes()
	if err != nil {
		return nil, fmt.Errorf("getting all time values: %w", err)
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

// getDateNowString returns today's date in the format of YYYY-MM-DD.
func GetDateNowString() string {
	y, m, d := time.Now().Date()

	currentDate := fmt.Sprintf("%d-", y)

	if m < 10 {
		currentDate += fmt.Sprintf("0%d-", m)
	} else {
		currentDate += fmt.Sprintf("%d-", m)
	}

	if d < 10 {
		currentDate += fmt.Sprintf("0%d", d)
	} else {
		currentDate += fmt.Sprintf("%d", d)
	}

	return currentDate
}
