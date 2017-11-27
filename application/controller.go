package application

import (
	"github.com/gin-gonic/gin"
)

type LogBookApplication struct {
	r *inbox
	d *messageDispatcher
}

func (app *LogBookApplication) AddApplicationToEngine(engine *gin.Engine) {
	app.r = &inbox{}
	app.d = &messageDispatcher{}

	// Create channel between inbox and messageDispatcher
	receiverToDispatcher := make(chan Message, 20)
	app.r.chanelToMessageDispatcher = receiverToDispatcher
	app.d.incoming = receiverToDispatcher

	// create messageDispatcher channel list
	app.d.channels = make(map[string]chan Message)
	go app.d.dispatch()

	engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		app.r.submit(context, logBookId)
	})

	engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		ForceCookie(context, logBookId)

		outbound := make(chan Message, 20)
		app.d.channels[logBookId] = outbound
		defer func() {
			delete(app.d.channels, logBookId)
		}()

		lb, _ := createLogBook(context.Writer, context.Request, outbound)
		lb.listen()
	})
}
