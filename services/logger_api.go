package services

import (
	"github.com/alexgunkel/logbook/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logContext interface {
	BindJSON(obj interface{}) error
}

func Log(c *gin.Context, toDispatcher chan<- entities.LogEvent) (err error)  {
	defer logEvent(c)
	e := &entities.LogEvent{}
	if err = c.BindJSON(e); nil != err {
		return
	}

	toDispatcher <- *e
	return
}

func logEvent(c *gin.Context)  {
	recover()
	c.Status(http.StatusBadRequest)
}