package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
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

	fmt.Printf("cleanup %s\n", id)
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
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var accessToken string
	var accessTokenSecret string

	err = stmtOut.QueryRow(id).Scan(&accessToken, &accessTokenSecret)

	fmt.Printf("id %s\n", id)
	fmt.Printf("accessToken %s\n", accessToken)
	fmt.Printf("accessTokenSecret %s\n", accessTokenSecret)

	ch := make(chan bool)
	d[id] = ch
	go connectStream(ch, id, accessToken, accessTokenSecret)
}

func connectStream(ch <-chan bool, id string, accessToken string, accessTokenSecret string) {
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
				fmt.Println("status: " + data.Text)
			case anaconda.StatusDeletionNotice:
				// pass
			case anaconda.EventTweet:
				// fmt.Println(data.Event.Event + " " + data.TargetObject.Text)
				bytes, _ := json.Marshal(data)
				fmt.Println("event: " + string(bytes))
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
