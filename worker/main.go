package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/anachronistic/apns"
	"github.com/benmanns/goworker"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func init() {
	goworker.Register("NotificationTweet", NotificationTweet)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}

func NotificationTweet(queue string, args ...interface{}) error {
	userIdStr, ok := args[0].(string)
	if !ok {
		fmt.Printf("invalid data:%v\n", args)
		return nil
	}

	screenName, ok := args[1].(string)
	if !ok {
		fmt.Printf("invalid data:%v\n", args)
		return nil
	}

	event, ok := args[2].(string)
	if !ok {
		fmt.Printf("invalid data:%v\n", args)
		return nil
	}

	text, ok := args[3].(string)
	if !ok {
		fmt.Printf("invalid data:%v\n", args)
		return nil
	}

	switch event {
	case "retweet":
		sendNotification(userIdStr, "@"+screenName+" さんがリツイート\n"+text)
	case "reply":
		sendNotification(userIdStr, "@"+screenName+" さんがリプライ\n"+text)
	default:
		sendNotification(userIdStr, "@"+screenName+" さんが"+event+"\n"+text)
	}

	return nil
}

func sendNotification(userIdStr string, message string) error {
	fmt.Printf("userId:%s message:%s\n", userIdStr, message)

	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE") // ex. justaway@tcp(192.168.0.10:3306)/justaway
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil
	}

	stmt := sq.
		Select("name, token, platform").
		From("notification_device").
		Where(sq.Eq{"user_id": userIdStr}).
		OrderBy("id DESC").
		Limit(100)

	sql, args, err := stmt.ToSql()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var token string
		var platform string
		err = rows.Scan(&name, &token, &platform)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		switch platform {
		case "APNS_SANDBOX":
			sendApns(token, message)
		}
	}

	return nil
}

func sendApns(token string, message string) {
	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = 1
	payload.Sound = "bingbong.aiff"

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	certificateFile := os.Getenv("JUSTAWAY_APNS_SANDBOX_CERT_PATH")  // apns-dev-cert.pem
	keyFile := os.Getenv("JUSTAWAY_APNS_SANDBOX_KEY_NOENC_PEM_PATH") // apns-dev-key-noenc.pem

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", certificateFile, keyFile)
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Token:", pn.DeviceToken)
	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)
}
