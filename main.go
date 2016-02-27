package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	fin := make(chan bool)

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		go connect(fin)
		c.JSON(200, gin.H{
			"message": "start streaming",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		fin <- true
		c.JSON(200, gin.H{
			"message": "stop streaming",
		})
	})
	r.Run()
}

func connect(fin <-chan bool) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	twitterStream := api.UserStream(nil)
	fmt.Println("connect")
	for {
		select {
		case x := <-twitterStream.C:
			switch data := x.(type) {
			case anaconda.FriendsList:
				// pass
			case anaconda.Tweet:
				fmt.Println(data.Text)
				fmt.Println("-----------")
			case anaconda.StatusDeletionNotice:
				// pass
			case anaconda.EventTweet:
				fmt.Println(data.Event.Event)
				fmt.Println(data.TargetObject.Text)
				fmt.Println("-----------")
				// pass
			default:
				fmt.Printf("unknown type(%T) : %v \n", x, x)
			}
		case <-fin:
			twitterStream.Stop()
			fmt.Println("fin")
			return
		}
	}
}
