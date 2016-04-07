package crawler

import (
	"database/sql"
	"fmt"
	"sync"
)
import _ "github.com/go-sql-driver/mysql"

const crawlerId = "1"
const dbSource = "justaway@tcp(192.168.0.10:3306)/justaway"

var m = new(sync.Mutex)
var d = map[string]chan bool{}

func Count() int {
	m.Lock()
	defer m.Unlock()

	count := 0
	for k, _ := range d {
		fmt.Printf("[Count] %s\n", k)
		count++
	}

	return count
}

func DisconnectAll() {
	m.Lock()
	defer m.Unlock()

	fmt.Println("[DisconnectAll] begin")

	for k, ch := range d {
		fmt.Printf("[DisconnectAll] %s\n", k)
		ch <- true
	}

	fmt.Println("[DisconnectAll] end")
}

func ConnectAll() {
	m.Lock()
	defer m.Unlock()

	fmt.Println("[ConnectAll] begin")

	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT user_id, access_token, access_token_secret FROM account WHERE crawler_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(crawlerId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var userId string
		var accessToken string
		var accessTokenSecret string
		err = rows.Scan(&userId, &accessToken, &accessTokenSecret)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("[ConnectAll] userId:%s accessToken:%s accessTokenSecret:%s\n", userId, accessToken, accessTokenSecret)

		ch := make(chan bool)
		d[userId] = ch
		go connectStream(ch, userId, accessToken, accessTokenSecret)
	}

	fmt.Println("[ConnectAll] end")
}

func Connect(userId string) {
	m.Lock()
	defer m.Unlock()

	_, ok := d[userId]
	if ok {
		return
	}

	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT access_token, access_token_secret FROM account WHERE user_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtOut.Close()

	var accessToken string
	var accessTokenSecret string

	err = stmtOut.QueryRow(userId).Scan(&accessToken, &accessTokenSecret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("[%s] accessToken:%s accessTokenSecret:%s\n", userId, accessToken, accessTokenSecret)

	ch := make(chan bool)
	d[userId] = ch
	go connectStream(ch, userId, accessToken, accessTokenSecret)
}

func Disconnect(userId string) {
	m.Lock()
	defer m.Unlock()

	ch, ok := d[userId]
	if ok {
		ch <- true
	}
}

func cleanup(userId string) {
	m.Lock()
	defer m.Unlock()

	ch, ok := d[userId]
	if ok {
		close(ch)
		delete(d, userId)
	}

	fmt.Printf("[%s] cleanup\n", userId)
}
