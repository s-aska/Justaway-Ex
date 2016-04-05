package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo"
	"io"
	"os"
)
import "text/template"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
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
	e.Get("/", index)
	e.Get("/count", count)
	e.Get("/signin/", signin)
	e.Get("/signin/callback", callback)
	e.Get("/api/activity/list.json", activity)
	e.Run("127.0.0.1:8002")
}
