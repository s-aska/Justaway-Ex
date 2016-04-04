package main

import (
	"github.com/labstack/echo"
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
	e.Get("/start", start)
	e.Get("/stop", stop)
	e.Run("127.0.0.1:8001")
}
