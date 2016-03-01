package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	e := echo.New()
	e.Debug()
	e.Get("/", index)
	e.Get("/signin", signin)
	e.Get("/callback", callback)
	e.Run("127.0.0.1:8002")
}
