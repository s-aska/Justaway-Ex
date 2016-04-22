package models

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

type Model struct {
	dbSource string
}

func New(dbSource string) *Model {
	return &Model{
		dbSource: dbSource,
	}
}

func (m *Model) CreateActivity(event string, targetId string, sourceId string, targetObjectId string, timestamp int64) {
	db, err := sql.Open("mysql", m.dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtIns, err := db.Prepare(`
		INSERT IGNORE INTO activity(
			event,
			target_id,
			source_id,
			target_object_id,
			created_at
		) VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(event, targetId, sourceId, targetObjectId, timestamp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if id == 0 {
		return
	}
}

func (m *Model) CreateRetweetActivity(event string, targetId string, sourceId string, targetObjectId string, retweeetedStatusId string, timestamp int64) {
	db, err := sql.Open("mysql", m.dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtIns, err := db.Prepare(`
		INSERT IGNORE INTO activity(
			event,
			target_id,
			source_id,
			target_object_id,
			retweeted_status_id,
			created_at
		) VALUES(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(event, targetId, sourceId, targetObjectId, retweeetedStatusId, timestamp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if id == 0 {
		return
	}
}

func (m *Model) DeleteTweetActivity(statusId string) {
	fmt.Printf("DeleteTweetActivity: %s\n", statusId)
	db, err := sql.Open("mysql", m.dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	func() {
		stmtDel, err := db.Prepare("DELETE FROM activity WHERE retweeted_status_id = ?")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer stmtDel.Close()

		_, err = stmtDel.Exec(statusId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	func() {
		stmtDel, err := db.Prepare("DELETE FROM activity WHERE target_object_id = ? AND event IN ('reply', 'retweet', 'quoted_tweet', 'favorite', 'favorited_retweet', 'retweeted_retweet')")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer stmtDel.Close()

		_, err = stmtDel.Exec(statusId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()
}

func (m *Model) DeleteFavoriteActivity(sourceId string, statusId string) {
	db, err := sql.Open("mysql", m.dbSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE event = 'favorite' AND source_id = ? AND target_object_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(sourceId, statusId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
