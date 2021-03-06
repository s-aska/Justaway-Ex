package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/s-aska/Justaway-Ex/web/routes"
	"io"
	"os"
	"strings"
	"text/template"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	consumerKey := os.Getenv("JUSTAWAY_EX_CONSUMER_KEY")
	consumerSecret := os.Getenv("JUSTAWAY_EX_CONSUMER_SECRET")
	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE") // ex. justaway@tcp(192.168.0.10:3306)/justaway
	callback := os.Getenv("JUSTAWAY_EX_CALLBACK")  // ex. https://justaway.info/signin/callback

	errors := []string{}
	if consumerKey == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CONSUMER_KEY=''")
	}
	if consumerSecret == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CONSUMER_SECRET=''")
	}
	if dbSource == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_DB_SOURCE=''")
	}
	if callback == "" {
		errors = append(errors, "$ export JUSTAWAY_EX_CALLBACK=''")
	}
	if len(errors) > 0 {
		fmt.Println(strings.Join(errors, "\n"))
		os.Exit(1)
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.SetRenderer(t)

	r := routes.New(dbSource, callback)

	e.Debug()
	e.Get("/signin/", r.SignIn)
	e.Get("/signin/callback", r.SignInCallback)
	e.Post("/api/revoke.json", r.Revoke)
	e.Put("/api/devices.json", r.ApiDeviceTokenRegister)
	e.Get("/api/activity/list.json", r.ApiActivityList)
	e.Get("/api/statuses/favoriters/ids.json", r.ApiFavoritersIds)
	e.Run(standard.New("127.0.0.1:8002"))
}
