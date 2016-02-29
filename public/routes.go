package main

import (
	"fmt"
	"net/http"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/garyburd/go-oauth/oauth"
)

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
		Token: token.(string),
		Secret: secret.(string),
	}
	cred, _, err := anaconda.GetCredentials(&tempCred, c.Query("oauth_verifier"))
	if err != nil {
		return c.String(200, err.Error())
	}

	session.Values["access_token"] = cred.Token
	session.Values["access_secret"] = cred.Secret
	session.Save(c.Request(), c.Response().Writer())

	return c.JSON(200, cred)
}
