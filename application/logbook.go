package application

import (
	"github.com/gin-gonic/gin"
)

type LogBookApplication struct {
	r *receiver
	d *dispatcher
}

func (app *LogBookApplication) AddApplicationToEngine(engine *gin.Engine) {
	app.r = &receiver{}
	app.d = &dispatcher{}

	// Create channel between receiver and dispatcher
	receiverToDispatcher := make(chan PostMessage, 20)
	app.r.cToDispatcher = receiverToDispatcher
	app.d.incoming = receiverToDispatcher

	// create dispatcher channel list
	app.d.channels = make(map[string]chan PostMessage)
	go app.d.dispatch()

	engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		app.r.Log(context, logBookId)
	})

	engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		outbound := make(chan PostMessage, 20)
		app.d.channels[logBookId] = outbound
		//ForceCookie(context)
		WebsocketHandler(context.Writer, context.Request, outbound)
	})
}
