package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	gsession "github.com/gorilla/sessions"
	"github.com/gin-gonic/contrib/sessions"
)

type MemcachedStore interface {
	sessions.Store
}

// client: memcache.Client.
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewMemcacheStore(client *memcache.Client, keyPrefix string, keyPairs ...[]byte) (MemcachedStore) {
	store := gsm.NewMemcacheStore(client, keyPrefix, keyPairs...)
	return &memcacheStore{store}
}

type memcacheStore struct {
	*gsm.MemcacheStore
}

func (c *memcacheStore) Options(options sessions.Options) {
	c.MemcacheStore.Options = &gsession.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}