package services

import (
	"github.com/alexgunkel/logbook/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Log(c *gin.Context, toDispatcher chan<- entities.PostMessage) (err error)  {
	m := &entities.PostMessage{}
	e := &entities.LogEvent{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}
	m.Event = *e

	h := entities.LogHeader{}
	h.Application = c.GetHeader(entities.LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(entities.LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(entities.LogHeaderRequestUri)

	m.Header = h

	toDispatcher <- *m
	return
}
