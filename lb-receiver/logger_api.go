package lb_receiver

import (
	"github.com/alexgunkel/logbook/lb-entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHeaderPrefix        string = "LogBook"
	LogHeaderAppIdentifier string = LogHeaderPrefix + "-App-Identifier"
	LogHeaderLoggerName    string = LogHeaderPrefix + "-Logger-Name"
	LogHeaderRequestUri    string = LogHeaderPrefix + "-Request-URI"
)

func Log(c *gin.Context, toDispatcher chan<- lb_entities.PostMessage) (err error)  {
	m := &lb_entities.PostMessage{}
	e := &lb_entities.LogEvent{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}
	m.Event = *e

	h := lb_entities.LogHeader{}
	h.Application = c.GetHeader(LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(LogHeaderRequestUri)

	m.Header = h

	toDispatcher <- *m
	return
}
