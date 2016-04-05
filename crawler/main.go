package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/s-aska/anaconda"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	e := echo.New()
	e.Debug()
	e.Get("/:id/start", start)
	e.Get("/:id/stop", stop)
	e.Run(standard.New("127.0.0.1:8001"))
}
