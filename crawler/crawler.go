package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/s-aska/anaconda"
	"strings"
	"sync"
	"time"
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
				// fmt.Printf("[%s] status: @%s `%s`\n", id, data.User.ScreenName, data.Text)
				if data.RetweetedStatus != nil && data.RetweetedStatus.User.IdStr == id {
					go createActivityWithReferenceId(
						id,
						data.RetweetedStatus.IdStr,
						"retweet",
						data.User.IdStr,
						data.IdStr,
						encodeJson(data))
				}

			case anaconda.DirectMessage:
				// fmt.Printf("[%s] message: @%s => @%s `%s`\n", id, data.SenderScreenName, data.RecipientScreenName, data.Text)

			case anaconda.StatusDeletionNotice:
				// fmt.Printf("[%s] status delete: %s:%s\n", id, data.UserIdStr, data.IdStr)
				go deleteActivityByStatusId(data.IdStr)
				go deleteActivityByReferenceId(data.IdStr)

			case anaconda.DirectMessageDeletionNotice:
				// fmt.Printf("[%s] message delete: %s:%s\n", id, data.UserId, data.IdStr)

			case anaconda.EventTweet:

				if data.Event.Event == "quoted_tweet" && data.TargetObject.QuotedStatus.User.IdStr == id {
					go createActivityWithReferenceId(
						id,
						data.TargetObject.QuotedStatus.IdStr,
						data.Event.Event,
						data.Event.Source.IdStr,
						data.TargetObject.IdStr,
						encodeJson(data))

				} else if data.TargetObject.User.IdStr == id {

					if data.Event.Event == "favorite" || data.Event.Event == "favorited_retweet" || data.Event.Event == "retweeted_retweet" {
						go createReactionActivity(
							id,
							data.TargetObject.IdStr,
							data.Event.Event,
							data.Event.Source.IdStr,
							encodeJson(data))

					} else if data.Event.Event == "unfavorite" {
						go deleteReactionActivity(
							data.TargetObject.IdStr,
							data.Event.Event,
							data.Event.Source.IdStr)
					}
				}

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

func createActivity(userId string, uniqueId string, data string) {
	fmt.Printf("[%s] createActivity: %s\n", userId, uniqueId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare(`
		INSERT IGNORE INTO activity(
			user_id,
			unique_id,
			data,
			created_on
		) VALUES(?, ?, ?, ?)
	`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userId, uniqueId, data, time.Now().Unix())
	if err != nil {
		panic(err.Error())
	}
}

func createReactionActivity(userId string, statusId string, event string, sourceUserId string, data string) {
	createActivity(userId, strings.Join([]string{"status", statusId, event, sourceUserId}, ":"), data)
}

func deleteReactionActivity(statusId string, event string, sourceUserId string) {
	deleteActivity(strings.Join([]string{"status", statusId, event, sourceUserId}, ":"))
}

func deleteActivity(uniqueId string) {
	fmt.Printf("[%s] deleteActivity: %s\n", "-", uniqueId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	ope := "="
	if strings.HasSuffix(uniqueId, "%") {
		ope = "LIKE"
	}

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE unique_id " + ope + " ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(uniqueId)
	if err != nil {
		panic(err.Error())
	}
}

func deleteActivityByStatusId(statusId string) {
	deleteActivity(strings.Join([]string{"status", statusId, "%"}, ":"))
}

func createActivityWithReferenceId(userId string, statusId string, event string, sourceUserId string, referenceId string, data string) {
	uniqueId := strings.Join([]string{"status", statusId, event, sourceUserId}, ":")
	fmt.Printf("[%s] createActivityWithReferenceId: %s %s\n", userId, uniqueId, referenceId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare(`
		INSERT IGNORE INTO activity(
			user_id,
			unique_id,
			reference_id,
			data,
			created_on
		) VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userId, uniqueId, referenceId, data, time.Now().Unix())
	if err != nil {
		panic(err.Error())
	}
}

func deleteActivityByReferenceId(referenceId string) {
	fmt.Printf("[%s] deleteActivityByReferenceId: %s\n", "-", referenceId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE reference_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(referenceId)
	if err != nil {
		panic(err.Error())
	}
}
