package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/gorilla/sessions"
)

const sessionName = "justaway_session"

func NewSessionStore(secure bool) *gsm.MemcacheStore {
	var memcacheClient = memcache.New("localhost:11211")
	var store = gsm.NewMemcacheStore(memcacheClient, "session_prefix_", []byte("secret-key-goes-here"))
	store.Options = &sessions.Options{
		MaxAge:   0,
		Path:     "/signin/",
		Secure:   secure,
		HttpOnly: true,
	}
	store.StoreMethod = gsm.StoreMethodGob
	return store
}
