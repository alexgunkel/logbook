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
	chanelToMessageDispatcher chan IncomingMessage
}

func (r *inbox) submit(c *gin.Context, logBookId string) (err error)  {
	m := createNewLogMessage(logBookId)
	e := &LogMessageBody{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}
	m.Body = *e

	h := createHeaderDataObjectFromHeaderData(c)

	m.Origin = h

	r.chanelToMessageDispatcher <- *m
	return
}

func createHeaderDataObjectFromHeaderData(c *gin.Context) (h HeaderData) {
	h = HeaderData{}
	h.Application = c.GetHeader(LogHeaderAppIdentifier)
	h.LoggerName = c.GetHeader(LogHeaderLoggerName)
	h.RequestUri = c.GetHeader(LogHeaderRequestUri)

	return
}
