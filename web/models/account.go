package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func (m *Model) CreateAccount(userIdStr string, name string, screenName string, accessToken string, accessTokenSecret string) string {
	now := time.Now()

	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	// fetch new crawler id
	var crawlerId string
	err = tx.QueryRow("SELECT id FROM crawler WHERE status = 'ACTIVE' ORDER BY id DESC LIMIT 1").Scan(&crawlerId)
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	_, err = tx.Exec(`
		INSERT INTO account(
			crawler_id,
			user_id,
			name,
			screen_name,
			access_token,
			access_token_secret,
			status,
			created_at,
			updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			screen_name = VALUES(screen_name),
			access_token = VALUES(access_token),
			access_token_secret = VALUES(access_token_secret),
			status = VALUES(status),
			updated_at = ?
	`, crawlerId, userIdStr, name, screenName, accessToken, accessTokenSecret, "ACTIVE", now.Unix(), 0, now.Unix())
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	apiToken := userIdStr + "-" + makeToken()

	_, err = tx.Exec(`
		INSERT INTO api_token(
			user_id,
			api_token,
			created_at,
			authenticated_at
		) VALUES(?, ?, ?, ?)
	`, userIdStr, apiToken, now.Unix(), now.Unix())
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	var crawlerUrl string
	err = tx.QueryRow(`
		SELECT crawler.url
		FROM account
		LEFT OUTER JOIN crawler ON crawler.id = account.crawler_id
		WHERE user_id = ? LIMIT 1
	`, userIdStr).Scan(&crawlerUrl)
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	tx.Commit()

	req, _ := http.NewRequest("GET", crawlerUrl+userIdStr+"/start", nil)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	fmt.Printf("start streaming user_id:%s screen_name:%s status:%s\n", userIdStr, screenName, res.Status)

	return apiToken
}

func makeToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
