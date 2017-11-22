package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/alexgunkel/logbook/entities"
)

type Webapp struct {
	engine *gin.Engine
}

func (app *Webapp) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	app.engine.ServeHTTP(w, req)
}

func Default() *gin.Engine {
	app := new(Webapp)
	app.engine = gin.Default()
	app.engine.GET("/logbook", InitLogBookClientApplication)
	app.engine.GET("/logbook/:client/logs", DisplayLogs)
	app.engine.POST("/logbook/:client/logs", Log)

	return app.engine
}

func DisplayLogs(c *gin.Context)  {
	_, err := c.Cookie("logbook")
	if nil != err {
		c.Redirect(http.StatusTemporaryRedirect, "../../logbook")
	}
}

func InitLogBookClientApplication(c *gin.Context)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		c.SetCookie("logbook", "asd", 0, "", "", false, false)
		identifier = "asd"
	}

	location := "logbook/" + identifier + "/logs"

	c.Redirect(http.StatusTemporaryRedirect, location)
}

func Log(c *gin.Context)  {
	e := &entities.LogEvent{}
	err := c.Bind(e)
	if nil != err {
		return
	}
}