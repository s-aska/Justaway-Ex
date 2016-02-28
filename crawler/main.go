package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	m := new(sync.Mutex)
	d := map[string]chan bool{}

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {

		m.Lock()
		defer m.Unlock()

		id := c.Param("id")
		_, ok := d[id]
		if ok {
			c.JSON(200, gin.H{
				"message": "exists streaming",
			})
		} else {
			d[id] = make(chan bool)
			go connect(d[id])
			c.JSON(200, gin.H{
				"message": "start streaming",
			})
		}
	})

	r.GET("/stop", func(c *gin.Context) {

		m.Lock()
		defer m.Unlock()

		id := c.Param("id")
		if fin, ok := d[id]; ok {
			fin <- true
			close(fin)
			delete(d, id)
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
			fmt.Println("disconnect")
			return
		}
	}
}
