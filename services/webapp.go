package services

import (
	"net/http"
	"github.com/alexgunkel/logbook/entities"
)

type webContext interface {
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Redirect(code int, location string)
}

type logContext interface {
	Bind(obj interface{}) error
}

func DisplayLogs(c webContext)  {
	_, err := c.Cookie("logbook")
	if nil != err {
		c.Redirect(http.StatusTemporaryRedirect, "../../logbook")
	}
}

func InitLogBookClientApplication(c webContext)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		c.SetCookie("logbook", "asd", 0, "", "", false, false)
		identifier = "asd"
	}

	location := "logbook/" + identifier + "/logs"

	c.Redirect(http.StatusTemporaryRedirect, location)
}

func Log(c logContext)  {
	e := &entities.LogEvent{}
	err := c.Bind(e)
	if nil != err {
		return
	}
}