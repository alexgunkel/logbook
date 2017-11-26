package application

import (
	"github.com/gin-gonic/gin"
)

func AddDispatcher(engine *gin.Engine)  {
	incoming := make(chan PostMessage, 20)
	//outbound := make(chan lb_entities.PostMessage, 10)

	engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		Log(context, incoming)
	})

	engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		ForceCookie(context)
		WebsocketHandler(context.Writer, context.Request, incoming)
	})
}
