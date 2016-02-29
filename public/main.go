package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	store := NewMemcacheStore(memcache.New("localhost:11211"), "justaway_session", []byte(consumerSecret))

	store.Options(sessions.Options{
		MaxAge: 0,
		Path: "/",
		Secure: true,
		HttpOnly: true,
	})

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("resources/*.templ.html")
	r.GET("/", index)
	// r.GET("/signin", signin)
	// r.GET("/callback", callback)
	r.Run(":10080")
}
