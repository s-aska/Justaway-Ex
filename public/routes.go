package main

import (
	"database/sql"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
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
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}
	store.StoreMethod = gsm.StoreMethodGob
	return store
}()

const session_name = "justaway_session"

func index(c *echo.Context) error {
	session, _ := store.Get(c.Request(), session_name)
	var count int
	value := session.Values["count"]
	if value == nil {
		count = 0
	} else {
		count = value.(int)
	}
	count = count + 1
	session.Values["count"] = count
	session.Save(c.Request(), c.Response().Writer())
	return c.Render(http.StatusOK, "index", "World")
}

func count(c *echo.Context) error {
	session, _ := store.Get(c.Request(), session_name)
	var count int
	value := session.Values["count"]
	if value == nil {
		count = 0
	} else {
		count = value.(int)
	}
	count = count + 1
	session.Values["count"] = count
	session.Save(c.Request(), c.Response().Writer())
	return c.String(200, fmt.Sprint(count))
}

func signin(c *echo.Context) error {
	url, tempCred, err := anaconda.AuthorizationURL("http://127.0.0.1:8002/callback")

	if err != nil {
		return c.String(200, err.Error())
	}

	session, _ := store.Get(c.Request(), session_name)
	session.Values["request_token"] = tempCred.Token
	session.Values["request_secret"] = tempCred.Secret
	session.Save(c.Request(), c.Response().Writer())

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func callback(c *echo.Context) error {
	session, _ := store.Get(c.Request(), session_name)
	token := session.Values["request_token"]
	secret := session.Values["request_secret"]
	tempCred := oauth.Credentials{
		Token:  token.(string),
		Secret: secret.(string),
	}
	cred, _, err := anaconda.GetCredentials(&tempCred, c.Query("oauth_verifier"))
	if err != nil {
		return c.String(200, err.Error())
	}

	session.Values["access_token"] = cred.Token
	session.Values["access_secret"] = cred.Secret
	session.Save(c.Request(), c.Response().Writer())

	api := anaconda.NewTwitterApi(cred.Token, cred.Secret)

	v := url.Values{}
	v.Set("include_entities", "false")
	v.Set("skip_status", "true")
	user, err := api.GetSelf(v)

	now := time.Now()
	fmt.Printf("id %s\n", user.Id)
	fmt.Printf("name %s\n", user.Name)
	fmt.Printf("ScreenName %s\n", user.ScreenName)

	db, err := sql.Open("mysql", "root:@/justaway")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO account(crawler_id, user_id, name, screen_name, access_token, access_token_secret, status, created_on, updated_on) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(1, user.Id, user.Name, user.ScreenName, cred.Token, cred.Secret, "ACTIVE", now.Unix(), 0)
	if err != nil {
		panic(err.Error())
	}

	return c.JSON(200, cred)
}
