package main

import (
	"crypto/tls"
	"net/http"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

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
	e.GET("/", func(c echo.Context) error {
		param := c.QueryParam("param")
		logger.Info(param)
		return c.String(http.StatusOK, ("param: " + param))
	})
	e.Start(":" + "5433")
}
