package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/s-aska/Justaway-Ex/crawler/handlers"
	"github.com/s-aska/anaconda"
)

func connectStream(ch <-chan bool, userId string, accessToken string, accessTokenSecret string) {
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)
	s := api.UserStream(nil)
	fmt.Printf("[%s] connecting...\n", userId)
	for {
		select {
		case x, ok := <-s.C:
			if !ok {
				fmt.Printf("[%s] disconnect\n", userId)
				s.Stop()
				cleanup(userId)
				return
			}
			switch data := x.(type) {
			case anaconda.FriendsList:
				fmt.Printf("[%s] connected\n", userId)
			case anaconda.Tweet:
				go handlers.HandlerTweet(userId, data)
			case anaconda.DirectMessage:
				go handlers.HandlerDirectMessage(userId, data)
			case anaconda.StatusDeletionNotice:
				go handlers.HandlerStatusDeletionNotice(data)
			case anaconda.DirectMessageDeletionNotice:
				go handlers.HandlerDirectMessageDeletionNotice(userId, data)
			case anaconda.EventTweet:
				go handlers.HandlerEventTweet(userId, data)
			case anaconda.EventList:
				fmt.Printf("[%s] eventList: %s %s\n", userId, data.Event.Event, encodeJson(data))
			case anaconda.Event:
				fmt.Printf("[%s] event: %s %s\n", userId, data.Event, encodeJson(data))
			case anaconda.DisconnectMessage:
				fmt.Printf("[%s] disconnectMessage\n", userId)
				s.Stop()
				cleanup(userId)
			default:
				fmt.Printf("[%s] unknown type(%T) : %v\n", userId, x, x)
			}
		case <-ch:
			fmt.Printf("[%s] stop\n", userId)
			s.Stop()
			cleanup(userId)
			return
		}
	}
}

func encodeJson(d interface{}) (j string) {
	b, _ := json.Marshal(d)
	return string(b)
}
