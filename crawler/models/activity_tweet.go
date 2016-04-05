package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

func makeTweetUniqueId(statusId string, event string, sourceUserId string) string {
	return strings.Join([]string{"tweet", statusId, event, sourceUserId}, ":")
}

func CreateTweetActivity(userId string, statusId string, event string, sourceUserId string, data string) {
	uniqueId := makeTweetUniqueId(statusId, event, sourceUserId)

	fmt.Printf("[%s] CreateTweetActivity: %s\n", userId, uniqueId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		fmt.Println(err.Error())
		return
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
		fmt.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userId, uniqueId, data, time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func DeleteTweetActivity(statusId string, event string, sourceUserId string) {
	uniqueId := makeTweetUniqueId(statusId, event, sourceUserId)

	fmt.Printf("[%s] DeleteTweetActivity: %s\n", "-", uniqueId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE unique_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(uniqueId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func DeleteTweetActivityByStatusId(statusId string) {
	uniqueId := strings.Join([]string{"tweet", statusId, "%"}, ":")

	fmt.Printf("[%s] DeleteTweetActivityByStatusId: %s\n", "-", uniqueId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE unique_id LIKE ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(uniqueId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func CreateTweetActivityWithReferenceId(userId string, statusId string, event string, sourceUserId string, referenceId string, data string) {
	uniqueId := makeTweetUniqueId(statusId, event, sourceUserId)

	fmt.Printf("[%s] CreateTweetActivityWithReferenceId: %s %s\n", userId, uniqueId, referenceId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		fmt.Println(err.Error())
		return
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
		fmt.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userId, uniqueId, referenceId, data, time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func DeleteTweetActivityByReferenceId(referenceId string) {
	fmt.Printf("[%s] DeleteTweetActivityByReferenceId: %s\n", "-", referenceId)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmtDel, err := db.Prepare("DELETE FROM activity WHERE reference_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(referenceId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
