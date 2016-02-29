package main

import (
	// "github.com/gorilla/sessions"
	// "github.com/bradfitz/gomemcache/memcache"
	// gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	// "github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	// "github.com/garyburd/go-oauth/oauth"
	// "net/http"
	// "fmt"
)

// var store = func() *gsm.MemcacheStore {
// 	var memcacheClient = memcache.New("localhost:11211")
// 	var store = gsm.NewMemcacheStore(memcacheClient, "session_prefix_", []byte("secret-key-goes-here"))
// 	store.Options = &sessions.Options{
// 		MaxAge: 0,
// 		Path: "/",
// 		Secure: false,
// 		HttpOnly: true,
// 	}
// 	store.StoreMethod = gsm.StoreMethodGob
// 	return store
// }()
// const session_name = "justaway_session"

func index(c *gin.Context) {
	session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count += 1
    }
    session.Set("count", count)
    session.Save()
	// session, _ := store.Get(c.Request, session_name)
	// var count int
	// value := session.Values["count"]
	// if value == nil {
	// 	count = 0
	// } else {
	// 	count = value.(int)
	// }
	// count = count + 1
	// session.Values["count"] = count
	// session.Save(c.Request, c.Writer)
	// c.JSON(200, gin.H{
	// 	"count": count,
	// })
	c.HTML(200, "index.templ.html", gin.H{
		"count": count,
	})
}

// func signin(c *gin.Context) {
// 	url, tempCred, err := anaconda.AuthorizationURL("http://127.0.0.1:10080/callback")

// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"message": err,
// 		})
// 		return
// 	}

// 	session, _ := store.Get(c.Request, session_name)
// 	session.Values["request_token"] = tempCred.Token
// 	session.Values["request_secret"] = tempCred.Secret
// 	session.Save(c.Request, c.Writer)

// 	c.Redirect(http.StatusMovedPermanently, url)
// }

// func callback(c *gin.Context) {
// 	session, _ := store.Get(c.Request, session_name)
// 	token := session.Values["request_token"]
// 	secret := session.Values["request_secret"]
// 	tempCred := oauth.Credentials{
// 		Token: token.(string),
// 		Secret: secret.(string),
// 	}
// 	cred, _, err := anaconda.GetCredentials(&tempCred, c.Param("oauth_verifier"))
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"message": err,
// 		})
// 		return
// 	}

// 	session.Values["access_token"] = cred.Token
// 	session.Values["access_secret"] = cred.Secret
// 	session.Save(c.Request, c.Writer)

// 	c.HTML(200, "index.templ.html", gin.H{})
// }
