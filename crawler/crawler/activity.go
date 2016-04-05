package crawler

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

func createActivity(userId string, statusId string, event string, sourceUserId string, data string) {
	uniqueId := strings.Join([]string{"status", statusId, event, sourceUserId}, ":")

	fmt.Printf("[%s] createActivity: %s\n", userId, uniqueId)

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

func deleteActivity(statusId string, event string, sourceUserId string) {
	uniqueId := strings.Join([]string{"status", statusId, event, sourceUserId}, ":")

	fmt.Printf("[%s] deleteActivity: %s\n", "-", uniqueId)

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

func deleteActivityByStatusId(statusId string) {
	uniqueId := strings.Join([]string{"status", statusId, "%"}, ":")

	fmt.Printf("[%s] deleteActivity: %s\n", "-", uniqueId)

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

func createActivityWithReferenceId(userId string, statusId string, event string, sourceUserId string, referenceId string, data string) {
	uniqueId := strings.Join([]string{"status", statusId, event, sourceUserId}, ":")

	fmt.Printf("[%s] createActivityWithReferenceId: %s %s\n", userId, uniqueId, referenceId)

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

func deleteActivityByReferenceId(referenceId string) {
	fmt.Printf("[%s] deleteActivityByReferenceId: %s\n", "-", referenceId)

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
