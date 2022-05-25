package main

import (
	"crypto/tls"
	"html/template"
	"io"
	"net/http"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	log := logrus.New()
	conn, err := tls.Dial("tcp", "996fdd9d-7a0e-49da-b84e-b9b7e055e6bb-ls.logit.io:10113", &tls.Config{RootCAs: nil})
	if err != nil {
		log.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "myappName"}))
	log.Hooks.Add(hook)
	logger := log.WithFields(logrus.Fields{
		"method": "main",
	})

	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		param := c.QueryParam("param")
		logger.Info(param)
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{
			"param": param,
		})
	})
	e.Start(":" + "5433")
}
