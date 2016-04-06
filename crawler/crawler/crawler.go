package crawler

import (
	"database/sql"
	"fmt"
	"sync"
)
import _ "github.com/go-sql-driver/mysql"

const dbSource = "justaway@tcp(192.168.0.10:3306)/justaway"

var m = new(sync.Mutex)
var d = map[string]chan bool{}

func Connect(id string) {
	m.Lock()
	defer m.Unlock()

	_, ok := d[id]
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

	err = stmtOut.QueryRow(id).Scan(&accessToken, &accessTokenSecret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("[%s] accessToken:%s accessTokenSecret:%s\n", id, accessToken, accessTokenSecret)

	ch := make(chan bool)
	d[id] = ch
	go connectStream(ch, id, accessToken, accessTokenSecret)
}

func Disconnect(id string) {
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
