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

	cleanup := func(id string) {
		m.Lock()
		defer m.Unlock()
		_, ok := d[id]
		if ok {
			close(d[id])
			delete(d, id)
		}
		fmt.Printf("cleanup %s\n", id)
	}

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
			return
		}

		fin := make(chan bool)
		d[id] = fin
		go connect(fin, id, cleanup)
		c.JSON(200, gin.H{
			"message": "start streaming",
		})
	})

	r.GET("/stop", func(c *gin.Context) {

		m.Lock()
		defer m.Unlock()

		id := c.Param("id")
		fin, ok := d[id]
		if !ok {
			c.JSON(200, gin.H{
				"message": "missing streaming",
			})
			return
		}

		fin <- true
		c.JSON(200, gin.H{
			"message": "stop streaming",
		})
	})
	r.Run()
}

func connect(fin <-chan bool, id string, cleanup func(id string)) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	stream := api.UserStream(nil)
	fmt.Println("connect")
	for {
		select {
		case x, ok := <-stream.C:
			if !ok {
				stream.Stop()
				fmt.Println("disconnect with error")
				cleanup(id)
				return
			}
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
			stream.Stop()
			fmt.Println("disconnect signal")
			cleanup(id)
			return
		}
	}
}
