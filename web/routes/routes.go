package routes

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/s-aska/Justaway-Ex/web/models"
	"strings"
)

type Router struct {
	dbSource     string
	callback     string
	sessionStore *gsm.MemcacheStore
}

func New(dbSource string, callback string) *Router {
	return &Router{
		dbSource:     dbSource,
		callback:     callback,
		sessionStore: NewSessionStore(strings.HasPrefix(callback, "https")),
	}
}

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

func (r *Router) loadSession(c echo.Context) (*sessions.Session, error) {
	return r.sessionStore.Get(c.Request().(*standard.Request).Request, sessionName)
}

func (r *Router) saveSession(c echo.Context, session *sessions.Session) {
	session.Save(c.Request().(*standard.Request).Request, c.Response().(*standard.Response).ResponseWriter)
}

func (r *Router) NewModel() *models.Model {
	return models.New(r.dbSource)
}
