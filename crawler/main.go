package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/s-aska/Justaway-Ex/crawler/crawler"
	"github.com/s-aska/Justaway-Ex/crawler/handlers"
	"github.com/s-aska/Justaway-Ex/crawler/models"
	"github.com/s-aska/anaconda"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	crawlerId := os.Getenv("JUSTAWAY_EX_CRAWLER_ID") // ex. 1
	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE")   // ex. justaway@tcp(192.168.0.10:3306)/justaway

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
