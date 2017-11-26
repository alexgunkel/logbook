package application

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHeaderPrefix        string = "LogBook"
	LogHeaderAppIdentifier string = LogHeaderPrefix + "-App-Identifier"
	LogHeaderLoggerName    string = LogHeaderPrefix + "-Logger-Name"
	LogHeaderRequestUri    string = LogHeaderPrefix + "-Request-URI"
)

func Log(c *gin.Context, toDispatcher chan<- PostMessage) (err error)  {
	m := &PostMessage{}
	e := &LogEvent{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}
	m.Event = *e

	h := LogHeader{}
	h.Application = c.GetHeader(LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(LogHeaderRequestUri)

	m.Header = h

	toDispatcher <- *m
	return
}
