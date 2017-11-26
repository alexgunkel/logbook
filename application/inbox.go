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

type inbox struct {
	chanelToMessageDispatcher chan LogMessage
}

func (r *inbox) submit(c *gin.Context, logBookId string) (err error)  {
	m := createNewLogMessage(logBookId)
	e := &LogEvent{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}
	m.Event = *e

	h := createOriginObjectFromHeaderData(c)

	m.Origin = h

	r.chanelToMessageDispatcher <- *m
	return
}

func createOriginObjectFromHeaderData(c *gin.Context) (h LogMessageOrigin) {
	h = LogMessageOrigin{}
	h.Application = c.GetHeader(LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(LogHeaderRequestUri)

	return
}
