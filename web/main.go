package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"io"
	"os"
)
import "text/template"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	dbSource := os.Getenv("JUSTAWAY_EX_DB_SOURCE") // ex. justaway@tcp(192.168.0.10:3306)/justaway
	callback := os.Getenv("JUSTAWAY_EX_CALLBACK") // ex. https://justaway.info/signin/callback

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.SetRenderer(t)

	r := NewRouter(dbSource, callback)

	e.Debug()
	e.Get("/signin/", r.signin)
	e.Get("/signin/callback", r.signinCallback)
	e.Get("/api/activity/list.json", r.activity)
	e.Run(standard.New("127.0.0.1:8002"))
}
