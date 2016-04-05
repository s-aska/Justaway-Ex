package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/s-aska/anaconda"
)

func connectStream(ch <-chan bool, id string, accessToken string, accessTokenSecret string) {
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)
	s := api.UserStream(nil)
	fmt.Printf("[%s] connect\n", id)
	for {
		select {
		case x, ok := <-s.C:
			if !ok {
				fmt.Printf("[%s] disconnect\n", id)
				s.Stop()
				cleanup(id)
				return
			}
			switch data := x.(type) {

			case anaconda.FriendsList:
				// pass

			case anaconda.Tweet:
				go handlerTweet(id, data)
			case anaconda.DirectMessage:
				go handlerDirectMessage(id, data)
			case anaconda.StatusDeletionNotice:
				go handlerStatusDeletionNotice(data)
			case anaconda.DirectMessageDeletionNotice:
				go handlerDirectMessageDeletionNotice(id, data)
			case anaconda.EventTweet:
				go handlerEventTweet(id, data)
			case anaconda.EventList:
				fmt.Printf("[%s] eventList: %s %s\n", id, data.Event.Event, encodeJson(data))
			case anaconda.Event:
				fmt.Printf("[%s] event: %s %s\n", id, data.Event, encodeJson(data))
			case anaconda.DisconnectMessage:
				fmt.Printf("[%s] disconnectMessage\n", id)
				s.Stop()
				cleanup(id)

			default:
				fmt.Printf("[%s] unknown type(%T) : %v\n", id, x, x)
			}
		case <-ch:
			fmt.Printf("[%s] stop\n", id)
			s.Stop()
			cleanup(id)
			return
		}
	}
}

func encodeJson(d interface{}) (j string) {
	b, _ := json.Marshal(d)
	return string(b)
}
