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
	e := &entities.LogEvent{}
	if err = c.BindJSON(e); nil != err {
		c.Status(http.StatusBadRequest)
		return
	}

	toDispatcher <- *e
	return
}
