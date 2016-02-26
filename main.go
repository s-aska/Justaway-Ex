package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	fmt.Println("begin")
	quit := make(chan bool)
	go func() {
		api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
		twitterStream := api.UserStream(nil)
		fmt.Println("connect")
		for {
			x := <-twitterStream.C
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
		}
		quit <- true
	}()
	<-quit
	fmt.Println("end")
}
