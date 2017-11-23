package services

import "github.com/alexgunkel/logbook/entities"

type logContext interface {
	Bind(obj interface{}) error
}

func Log(c logContext)  {
	e := &entities.LogEvent{}
	err := c.Bind(e)
	if nil != err {
		return
	}
}
