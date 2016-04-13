package routes

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

func (r *Router) Signin(c echo.Context) error {
	url, tempCred, err := anaconda.AuthorizationURL(r.callback)

	if err != nil {
		return c.String(200, err.Error())
	}

	session, _ := r.loadSession(c)
	session.Values["request_token"] = tempCred.Token
	session.Values["request_secret"] = tempCred.Secret
	r.saveSession(c, session)

	return c.Redirect(http.StatusTemporaryRedirect, url+"&force_login=true")
}

func (r *Router) SigninCallback(c echo.Context) error {
	session, _ := r.loadSession(c)
	token := session.Values["request_token"]
	secret := session.Values["request_secret"]
	tempCred := oauth.Credentials{
		Token:  token.(string),
		Secret: secret.(string),
	}
	session.Values["request_token"] = ""
	session.Values["request_secret"] = ""
	r.saveSession(c, session)

	cred, _, err := anaconda.GetCredentials(&tempCred, c.QueryParam("oauth_verifier"))
	if err != nil {
		return c.String(200, err.Error())
	}

	api := anaconda.NewTwitterApi(cred.Token, cred.Secret)

	v := url.Values{}
	v.Set("include_entities", "false")
	v.Set("skip_status", "true")
	user, err := api.GetSelf(v)

	fmt.Printf("signin success user_id:%s screen_name:%s name:%s\n", user.IdStr, user.ScreenName, user.Name)

	apiToken := r.NewModel().CreateAccount(user.IdStr, user.Name, user.ScreenName, cred.Token, cred.Secret)

	return c.Render(http.StatusOK, "index", apiToken)
}
