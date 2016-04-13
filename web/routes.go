package main

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	// "encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	sq "gopkg.in/Masterminds/squirrel.v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

const sessionName = "justaway_session"

type Router struct {
	dbSource     string
	callback     string
	sessionStore *gsm.MemcacheStore
}

// NullUint64 is a sql.Scanner for unsigned ints.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Scan implements the sql.Scanner interface.
func (n *NullUint64) Scan(src interface{}) error {
	if src == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	s := asString(src)
	var err error
	n.Uint64, err = strconv.ParseUint(s, 10, 64)
	return err
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	return fmt.Sprintf("%v", src)
}

type JsonNullUInt64 struct {
	NullUint64
}

func (v JsonNullUInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Uint64)
	} else {
		return json.Marshal(nil)
	}
}

type Activity struct {
	Id                uint64         `json:"id"`
	Event             string         `json:"event"`
	TargetId          uint64         `json:"target_id"`
	SourceId          uint64         `json:"source_id"`
	TargetObjectId    uint64         `json:"target_object_id"`
	RetweetedStatusId JsonNullUInt64 `json:"retweeted_status_id"`
	CreatedOn         int            `json:"created_on"`
}

func NewRouter(dbSource string, callback string) *Router {
	return &Router{
		dbSource:     dbSource,
		callback:     callback,
		sessionStore: NewSessionStore(strings.HasPrefix(callback, "https")),
	}
}

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

func (r *Router) signin(c echo.Context) error {
	url, tempCred, err := anaconda.AuthorizationURL(r.callback)

	if err != nil {
		return c.String(200, err.Error())
	}

	session, _ := r.sessionStore.Get(c.Request().(*standard.Request).Request, sessionName)
	session.Values["request_token"] = tempCred.Token
	session.Values["request_secret"] = tempCred.Secret
	session.Save(c.Request().(*standard.Request).Request, c.Response().(*standard.Response).ResponseWriter)

	return c.Redirect(http.StatusTemporaryRedirect, url+"&force_login=true")
}

func (r *Router) signinCallback(c echo.Context) error {
	session, _ := r.sessionStore.Get(c.Request().(*standard.Request).Request, sessionName)
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

	db, err := sql.Open("mysql", r.dbSource)
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

	apiToken := r.makeToken()

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

func (r *Router) makeToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (r *Router) activity(c echo.Context) error {
	apiToken := c.Request().Header().Get("X-Justaway-API-Token")
	if apiToken == "" {
		return c.String(401, "Missing X-Justaway-API-Token header")
	}

	db, err := sql.Open("mysql", r.dbSource)
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

	toId := func(id string) string {
		if id == "" {
			return ""
		} else if strings.Contains(id, ":") {
			fields := strings.Split(id, ":")
			stmtOut, err := db.Prepare("SELECT id FROM activity WHERE target_object_id = ? AND event = ? AND source_id = ? LIMIT 1")
			if err != nil {
				fmt.Println(err.Error())
				return ""
			}
			defer stmtOut.Close()
			var dbId string
			fmt.Printf("target_object_id:%s event:%s source_id:%s\n", fields[0], fields[1], fields[2])
			err = stmtOut.QueryRow(fields[0], fields[1], fields[2]).Scan(&dbId)
			if err != nil {
				fmt.Println(err.Error())
				return ""
			}
			return dbId
		} else {
			return id
		}
	}

	maxId := toId(c.QueryParam("max_id"))
	sinceId := toId(c.QueryParam("since_id"))

	stmt := sq.
		Select("id, event, target_id, source_id, target_object_id, retweeted_status_id, created_on").
		From("activity").
		Where(sq.Eq{"target_id": userIdStr}).
		OrderBy("id DESC").
		Limit(200)

	if maxId != "" {
		stmt = stmt.Where(sq.LtOrEq{"id": maxId})
	}

	if sinceId != "" {
		stmt = stmt.Where(sq.GtOrEq{"id": sinceId})
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("maxId:%s sinceId:%s sql:%s\n", maxId, sinceId, sql)

	rows, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	events := []*Activity{}

	for rows.Next() {
		var id uint64
		var event string
		var targetId uint64
		var sourceId uint64
		var targetObjectId uint64
		var retweetedStatusId JsonNullUInt64
		var createdOn int
		err = rows.Scan(&id, &event, &targetId, &sourceId, &targetObjectId, &retweetedStatusId, &createdOn)
		if err != nil {
			panic(err.Error())
		}
		a := &Activity{
			Id:                id,
			Event:             event,
			TargetId:          targetId,
			SourceId:          sourceId,
			TargetObjectId:    targetObjectId,
			RetweetedStatusId: retweetedStatusId,
			CreatedOn:         createdOn,
		}
		events = append(events, a)
	}

	return c.JSON(200, events)
}
