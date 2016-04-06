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

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.SetRenderer(t)

	e.Debug()
	e.Get("/signin/", signin)
	e.Get("/signin/callback", callback)
	e.Get("/api/activity/list.json", activity)
	e.Run(standard.New("127.0.0.1:8002"))
}
