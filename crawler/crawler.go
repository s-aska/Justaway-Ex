package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"sync"
)

var m = new(sync.Mutex)
var d = map[string]chan bool{}

func disconnect(id string) {
	m.Lock()
	defer m.Unlock()

	ch, ok := d[id]
	if ok {
		ch <- true
	}
}

func cleanup(id string) {
	m.Lock()
	defer m.Unlock()

	ch, ok := d[id]
	if ok {
		close(ch)
		delete(d, id)
	}

	fmt.Printf("cleanup %s\n", id)
}

func connect(id string) {
	m.Lock()
	defer m.Unlock()

	_, ok := d[id]
	if ok {
		return
	}

	ch := make(chan bool)
	d[id] = ch
	go connectStream(ch, id)
}

func connectStream(ch <-chan bool, id string) {
	// TODO: get by DB
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)
	s := api.UserStream(nil)
	fmt.Printf("connect %s\n", id)
	for {
		select {
		case x, ok := <-s.C:
			if !ok {
				fmt.Printf("disconnect %s\n", id)
				s.Stop()
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
		case <-ch:
			fmt.Printf("stop %s\n", id)
			s.Stop()
			cleanup(id)
			return
		}
	}
}
