package logbook

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

const CONTENT_TYPE_JSON  = "application/json"

type logEvent struct {
	m string
}

func SetUpRouter() *gin.Engine {
	engine := gin.Default()
	engine.GET("/display", DisplayLogBook)
	engine.POST("/logbook/:client/logs", Log)

	return engine
}

func DisplayLogBook(c *gin.Context)  {
	c.SetCookie("logbook", "asd", 0, "", "", false, false)
	c.Header("Content-Type", CONTENT_TYPE_JSON)
	c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
}

func Log(c *gin.Context)  {
	e := &logEvent{}
	err := c.Bind(e)
	if nil != err {
		return
	}
}
