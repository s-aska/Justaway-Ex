package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/benmanns/goworker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"os"
)

var apnsClientProduction *apns2.Client
var apnsClientDevelopment *apns2.Client

func init() {
	goworker.Register("NotificationTweet", NotificationTweet)

	pemFile := os.Getenv("JUSTAWAY_APNS_PEM_PATH") // apns.pem

	cert, pemErr := certificate.FromPemFile(pemFile, "")
	if pemErr != nil {
		fmt.Println("Cert Error:", pemErr)
	}

	apnsClientProduction = apns2.NewClient(cert).Production()
	apnsClientDevelopment = apns2.NewClient(cert).Development()
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

	message := makeNotificationMessage(event, screenName, text)
	if err := sendNotification(userIdStr, message); err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return nil
}

func makeNotificationMessage(event string, screenName string, text string) string {
	switch event {
	case "retweet", "retweeted_retweet":
		return "@" + screenName + "さんがリツイートしました\n" + text
	case "reply":
		return "@" + screenName + "さんからのリプライ\n" + text
	case "favorite", "favorited_retweet":
		return "@" + screenName + "さんがいいねしました\n" + text
	case "quoted_tweet":
		return "@" + screenName + "さんが引用しました\n" + text
	default:
		return "@" + screenName + "さんが" + event + "\n" + text
	}
}

func sendNotification(userIdStr string, message string) error {
	fmt.Printf("userId:%s message:%s\n", userIdStr, message)

	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE") // ex. justaway@tcp(192.168.0.10:3306)/justaway
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return err
	}

	stmt := sq.
		Select("name, token, platform").
		From("notification_device").
		Where(sq.Eq{"user_id": userIdStr}).
		OrderBy("id DESC").
		Limit(100)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return err
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
		case "APNS":
			sendApns(apnsClientProduction, token, message)
		case "APNS_SANDBOX":
			sendApns(apnsClientDevelopment, token, message)
		}
	}

	return nil
}

func sendApns(client *apns2.Client, token string, message string) {
	payload := payload.NewPayload().Alert(message)
	notification := &apns2.Notification{}
	notification.DeviceToken = token
	notification.Topic = "pw.aska.Justaway"
	notification.Payload = payload
	res, err := client.Push(notification)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("APNs ID:", res.ApnsID)
}
