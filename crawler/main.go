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

	chans := map[string]chan bool{}

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		id := "abc"
		if chans[id] {
			c.JSON(200, gin.H{
				"message": "exists streaming",
			})
		} else {
			chans[id] = make(chan bool)
			go connect(chans[id])
			c.JSON(200, gin.H{
				"message": "start streaming",
			})
		}
	})
	r.GET("/stop", func(c *gin.Context) {
		id := "abc"
		if fin, ok := chans[id]; ok {
			fin <- true
			delete(chans, id)
			c.JSON(200, gin.H{
				"message": "stop streaming",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "missing streaming",
			})
		}
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
