# Email News

## Features

-   Scrapes data from Hacker News, Morning Brew, and TLDR Newsletters
-   Basic web interface for scraped data
-   Daily emails with scraped data

## Running the program

Assuming that the program was merely cloned, we would need to...

1. Set up a local database file
2. Set up the .env
3. Install [go](https://go.dev/doc/install)
4. Generate a [templ](https://templ.guide/quick-start/installation) template

### Set up: Database

1. Download the `emne.db` file from https://densola.github.io/files/
2. Add it into in this repo's root
3. Rename it to `db.db`

### Set up: .env

Create a `.env` file inside the repo's root containing the following content:

Details:

```.env
PORT=           # Port for the program
TIME=           # 24-hour formatted UTC time for when the program should scrape/send out emails
MAILFROM=       # Email sender's email address
MAILPASS=       # Email sender's password
MAILTO=         # Email receiver's email address
MAILHOST=       # Email service provider host/server
MAILPORT=       # Email service provider port
```

Sample:

```.env
PORT=8080
TIME="23:00"
MAILFROM="mailfrom@email.com"
MAILPASS="12345"
MAILTO="mailto@email.com"
MAILHOST="mail.site.com"
MAILPORT="123"
```

### Set up: Go and templ

1. Install Go from https://go.dev/doc/install
2. Install templ from https://templ.guide/quick-start/installation
3. Open a terminal, then...
    1. cd into this repo's root
    2. Run `templ generate`
    3. Run `go run .`

After this, the program should be running locally at the specified port.

## Building and deploying the program

TODO...
