package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/s-aska/anaconda"
	// "os"
	"sync"
)
import _ "github.com/go-sql-driver/mysql"

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

	fmt.Printf("[%s] cleanup\n", id)
}

func connect(id string) {
	m.Lock()
	defer m.Unlock()

	_, ok := d[id]
	if ok {
		return
	}

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT access_token, access_token_secret FROM account WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var accessToken string
	var accessTokenSecret string

	err = stmtOut.QueryRow(id).Scan(&accessToken, &accessTokenSecret)

	fmt.Printf("[%s] accessToken:%s accessTokenSecret:%s\n", id, accessToken, accessTokenSecret)

	ch := make(chan bool)
	d[id] = ch
	go connectStream(ch, id, accessToken, accessTokenSecret)
}

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
				fmt.Printf("[%s] status: @%s `%s`\n", id, data.User.ScreenName, data.Text)
			case anaconda.DirectMessage:
				fmt.Printf("[%s] message: @%s => @%s `%s`\n", id, data.SenderScreenName, data.RecipientScreenName, data.Text)
			case anaconda.StatusDeletionNotice:
				fmt.Printf("[%s] status delete: %s:%s\n", id, data.UserIdStr, data.IdStr)
			case anaconda.DirectMessageDeletionNotice:
				fmt.Printf("[%s] message delete: %s:%s\n", id, data.UserId, data.IdStr)
				bytes, _ := json.Marshal(data)
				fmt.Printf("[%s] message delete: %s\n", id, string(bytes))
			case anaconda.EventTweet:
				bytes, _ := json.Marshal(data)
				fmt.Printf("[%s] eventTweet: %s %s\n", id, data.Event.Event, string(bytes))
			case anaconda.EventList:
				bytes, _ := json.Marshal(data)
				fmt.Printf("[%s] eventList: %s %s\n", id, data.Event.Event, string(bytes))
			case anaconda.Event:
				bytes, _ := json.Marshal(data)
				fmt.Printf("[%s] event: %s %s", id, data.Event, string(bytes))
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
