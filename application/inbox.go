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
	chanelToMessageDispatcher chan Message
}

func (r *inbox) submit(c *gin.Context, logBookId string) (err error)  {
	m := createNewLogMessage(logBookId)
	e := &Event{}
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

func createOriginObjectFromHeaderData(c *gin.Context) (h Origin) {
	h = Origin{}
	h.Application = c.GetHeader(LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(LogHeaderRequestUri)

	return
}
