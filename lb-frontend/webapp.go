package lb_frontend

import (
	"net/http"
	"strconv"
)

type webContext interface {
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Redirect(code int, location string)
}

type IdGenerator struct {
	lastIdentifier int64
}

func (app *IdGenerator) getNewIdentifier() string {
	app.lastIdentifier++
	return strconv.FormatInt(app.lastIdentifier, 10)
}

func DisplayLogs(c webContext, gen *IdGenerator)  {
	_, err := c.Cookie("logbook")
	if nil != err {
		identifier := gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
		c.Redirect(http.StatusTemporaryRedirect, "../../logbook/" + identifier + "/logs")
	}
}

func InitLogBookClientApplication(c webContext, gen *IdGenerator)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	location := "logbook/" + identifier + "/logs"

	c.Redirect(http.StatusTemporaryRedirect, location)
}