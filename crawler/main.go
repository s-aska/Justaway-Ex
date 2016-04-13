package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/s-aska/Justaway-Ex/crawler/crawler"
	"github.com/s-aska/Justaway-Ex/crawler/handlers"
	"github.com/s-aska/Justaway-Ex/crawler/models"
	"github.com/s-aska/anaconda"
	"os"
	"strings"
)

func main() {
	consumerKey := os.Getenv("JUSTAWAY_EX_CONSUMER_KEY")
	consumerSecret := os.Getenv("JUSTAWAY_EX_CONSUMER_SECRET")
	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE")   // ex. justaway@tcp(192.168.0.10:3306)/justaway
	crawlerId := os.Getenv("JUSTAWAY_EX_CRAWLER_ID") // ex. 1

	errors := []string{}
	if consumerKey == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CONSUMER_KEY=''")
	}
	if consumerSecret == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CONSUMER_SECRET=''")
	}
	if dbSource == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_DB_SOURCE=''")
	}
	if crawlerId == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CRAWLER_ID=''")
	}
	if len(errors) > 0 {
		fmt.Println(strings.Join(errors, "\n"))
		os.Exit(1)
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	m := models.New(dbSource)
	h := handlers.New(m)
	c := crawler.New(crawlerId, dbSource, h)
	r := &Router{crawler: c}

	e := echo.New()
	e.Debug()
	e.Get("/startup", r.startup)
	e.Get("/shutdown", r.shutdown)
	e.Get("/status", r.status)
	e.Get("/:id/start", r.start)
	e.Get("/:id/stop", r.stop)
	e.Run(standard.New("127.0.0.1:8001"))
}
