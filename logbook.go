package logbook

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

const CONTENT_TYPE_JSON  = "application/json"
const CONTENT_TYPE_HTML  = "text/html"

type logEvent struct {
	m string
}

func SetUpRouter() *gin.Engine {
	engine := gin.Default()
	engine.GET("/logbook", InitLogBookClientApplication)
	engine.GET("/logbook/:client/logs", DisplayLogs)
	engine.POST("/logbook/:client/logs", Log)

	return engine
}

func DisplayLogs(c *gin.Context)  {
	_, err := c.Cookie("logbook")
	if nil != err {
		c.Redirect(http.StatusTemporaryRedirect, "../../logbook")
	}
}

func InitLogBookClientApplication(c *gin.Context)  {
	c.SetCookie("logbook", "asd", 0, "", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, "asd/logs")
}

func Log(c *gin.Context)  {
	e := &logEvent{}
	err := c.Bind(e)
	if nil != err {
		return
	}
}
