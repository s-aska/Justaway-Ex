package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	r := gin.Default()
	r.GET("/start", start)
	r.GET("/stop", stop)
	r.Run()
}
