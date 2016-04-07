package main

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var store = func() *gsm.MemcacheStore {
	var memcacheClient = memcache.New("localhost:11211")
	var store = gsm.NewMemcacheStore(memcacheClient, "session_prefix_", []byte("secret-key-goes-here"))
	store.Options = &sessions.Options{
		MaxAge:   0,
		Path:     "/signin/",
		Secure:   true,
		HttpOnly: true,
	}
	store.StoreMethod = gsm.StoreMethodGob
	return store
}()

const sessionName = "justaway_session"
const dbSource = "justaway@tcp(192.168.0.10:3306)/justaway"

func signin(c echo.Context) error {
	url, tempCred, err := anaconda.AuthorizationURL("https://justaway.info/signin/callback")

	if err != nil {
		return c.String(200, err.Error())
	}

	session, _ := store.Get(c.Request().(*standard.Request).Request, sessionName)
	session.Values["request_token"] = tempCred.Token
	session.Values["request_secret"] = tempCred.Secret
	session.Save(c.Request().(*standard.Request).Request, c.Response().(*standard.Response).ResponseWriter)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func callback(c echo.Context) error {
	session, _ := store.Get(c.Request().(*standard.Request).Request, sessionName)
	token := session.Values["request_token"]
	secret := session.Values["request_secret"]
	tempCred := oauth.Credentials{
		Token:  token.(string),
		Secret: secret.(string),
	}
	cred, _, err := anaconda.GetCredentials(&tempCred, c.QueryParam("oauth_verifier"))
	if err != nil {
		return c.String(200, err.Error())
	}

	session.Values["access_token"] = cred.Token
	session.Values["access_secret"] = cred.Secret
	session.Save(c.Request().(*standard.Request).Request, c.Response().(*standard.Response).ResponseWriter)

	api := anaconda.NewTwitterApi(cred.Token, cred.Secret)

	v := url.Values{}
	v.Set("include_entities", "false")
	v.Set("skip_status", "true")
	user, err := api.GetSelf(v)

	now := time.Now()
	fmt.Printf("callback user_id:%s screen_name:%s name:%s\n", user.Id, user.ScreenName, user.Name)

	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare(`
		INSERT INTO account(
			crawler_id,
			user_id,
			name,
			screen_name,
			api_token,
			access_token,
			access_token_secret,
			status,
			created_on,
			updated_on
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			screen_name = VALUES(screen_name),
			api_token = VALUES(api_token),
			access_token = VALUES(access_token),
			access_token_secret = VALUES(access_token_secret),
			status = VALUES(status),
			updated_on = ?
	`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	apiToken := makeToken()

	_, err = stmtIns.Exec(1, user.Id, user.Name, user.ScreenName, apiToken, cred.Token, cred.Secret, "ACTIVE", now.Unix(), 0, now.Unix())
	if err != nil {
		panic(err.Error())
	}

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8001/"+user.IdStr+"/start", nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("request user_id:%s screen_name:%s res:%s\n", user.IdStr, user.ScreenName, string(byteArray))

	return c.Render(http.StatusOK, "index", user.IdStr+"-"+apiToken)
}

func makeToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func activity(c echo.Context) error {
	apiToken := c.Request().Header().Get("X-Justaway-API-Token")
	if apiToken == "" {
		return c.String(401, "Missing X-Justaway-API-Token header")
	}

	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT user_id FROM account WHERE api_token = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var userIdStr string

	err = stmtOut.QueryRow(apiToken).Scan(&userIdStr)
	if err != nil {
		return c.String(401, "Invalid X-Justaway-API-Token header")
	}

	stmtData, err := db.Prepare("SELECT id, data FROM activity WHERE user_id = ? ORDER BY id DESC LIMIT 200")
	if err != nil {
		panic(err.Error())
	}
	defer stmtData.Close()

	rows, err := stmtData.Query(userIdStr)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	maxIdStr := "null"
	minIdStr := "null"
	events := ""
	count := 0
	for rows.Next() {
		count++
		var id string
		var data string
		err = rows.Scan(&id, &data)
		if err != nil {
			panic(err.Error())
		}
		if events == "" {
			events = data
		} else {
			events = events + "," + data
		}
		if maxIdStr == "null" {
			maxIdStr = "\""+id+"\""
		}
		minIdStr = "\""+id+"\""
	}

	return c.String(200, "{\"events\":["+events+"],\"max_id_str\":"+maxIdStr+",\"min_id_str\":"+minIdStr+"}")
}
